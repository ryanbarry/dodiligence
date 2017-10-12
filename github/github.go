package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/jtacoma/uritemplates"
)

var c = &http.Client{Timeout: (10 * time.Second)}

type Github struct {
	username, token string
}

func NewClient(username, token string) Github {
	return Github{username, token}
}

func (gh Github) GetAllRepos(org string) (repos []GHRepo, err error) {
	req, err := http.NewRequest("GET", "https://api.github.com/orgs/"+org+"/repos", nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(gh.username, gh.token)

	n, _, repos, err := getRepoPage(req)
	for {
		//fmt.Printf("next: %s, last: %s\n", n, l)
		if n == "" {
			break
		}

		req, err = http.NewRequest("GET", n, nil)
		if err != nil {
			return
		}
		req.SetBasicAuth(gh.username, gh.token)

		var repoPage []GHRepo
		n, _, repoPage, err = getRepoPage(req)
		if err != nil {
			return
		}

		repos = append(repos, repoPage...)
	}

	return
}

var linkMatcher = regexp.MustCompile("<(https?://[a-zA-Z0-9/?=.-]+)>; rel=\"(next|last)\"(, )?")

func getRepoPage(req *http.Request) (next, last string, repos []GHRepo, err error) {
	res, err := c.Do(req)
	if err != nil {
		return
	}

	//fmt.Printf("Got %s\n", res.Status)

	lh := res.Header[http.CanonicalHeaderKey("link")]
	links := make(map[string]string)

	for _, l := range lh {
		sm := linkMatcher.FindAllStringSubmatch(l, -1)
		for _, m := range sm {
			// links["next"] = "http://api.github.com/orgs/metamx/repos?page=2"
			links[m[2]] = m[1]
		}
	}
	next = links["next"]
	last = links["last"]

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &repos)
	if err != nil {
		return
	}

	return
}

func (gh Github) DownloadRepoArchive(repo GHRepo, archiveFormat GHArchiveFormat, ref string) error {
	of, err := os.OpenFile(repo.Name+"_"+repo.DefaultBranch+archiveFormatFileEx[archiveFormat], os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	t, err := uritemplates.Parse(repo.ArchiveUrl)
	if err != nil {
		return err
	}
	v := map[string]interface{}{"archive_format": archiveFormatValues[archiveFormat]}
	if ref != "" {
		v["ref"] = ref
	}
	ex, err := t.Expand(v)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", ex, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(gh.username, gh.token)

	res, err := http.DefaultClient.Do(req) // use DefaultClient since it has no timeouts and downloads may take awhile
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Got %s!", res.Status)
	}

	_, err = io.Copy(of, res.Body)
	return err
}
