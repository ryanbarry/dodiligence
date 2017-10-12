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

func (gh Github) GetRepo(org, repoName string) (repo GHRepo, err error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+org+"/"+repoName, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(gh.username, gh.token)

	res, err := c.Do(req)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &repo)
	return
}

func (gh Github) GetAllRepos(org string) (repos []GHRepo, err error) {
	var req *http.Request
	next := "https://api.github.com/orgs/" + org + "/repos"
	i := 1 // GH paginates starting at 1
	repos = []GHRepo{}

	for {
		if next == "" {
			break
		}

		req, err = http.NewRequest("GET", next, nil)
		if err != nil {
			return
		}
		req.SetBasicAuth(gh.username, gh.token)

		var repoPage []GHRepo
		next, _, repoPage, err = getRepoPage(req)
		if err != nil {
			err = fmt.Errorf("Error getting page %d of repo list: %s", i, err.Error())
			return
		}

		repos = append(repos, repoPage...)
		i++
	}

	return
}

var linkMatcher = regexp.MustCompile("<(https?://[a-zA-Z0-9/?=.-]+)>; rel=\"(next|last)\"(, )?")

func getRepoPage(req *http.Request) (next, last string, repos []GHRepo, err error) {
	res, err := c.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		err = fmt.Errorf("Got %s", res.Status)
		return
	}

	lh := res.Header[http.CanonicalHeaderKey("link")]
	links := make(map[string]string)

	for _, l := range lh {
		sm := linkMatcher.FindAllStringSubmatch(l, -1)
		for _, m := range sm {
			// links["next"] = "http://api.github.com/orgs/MYCOOLORGANIZATION/repos?page=2"
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

	localFileName := repo.Name + "_" + repo.DefaultBranch + archiveFormatFileEx[archiveFormat]
	of, err := os.OpenFile(localFileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		of.Close()
		return err
	}

	req, err := http.NewRequest("GET", ex, nil)
	if err != nil {
		of.Close()
		os.Remove(localFileName)
		return err
	}
	req.SetBasicAuth(gh.username, gh.token)

	// use http.DefaultClient since it has no timeouts and downloads may take awhile
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		of.Close()
		os.Remove(localFileName)
		return err
	}

	if res.StatusCode != 200 {
		of.Close()
		os.Remove(localFileName)
		return fmt.Errorf("Got %s!", res.Status)
	}

	_, err = io.Copy(of, res.Body)
	of.Close()
	return err
}
