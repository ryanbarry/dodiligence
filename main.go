package main

import (
	"flag"
	"fmt"

	"github.com/ryanbarry/dodiligence/github"
)

func main() {
	var u = flag.String("username", "", "github username")
	var p = flag.String("token", "", "github personal auth token")
	flag.Parse()

	if len(*u) == 0 || len(*p) == 0 {
		panic(fmt.Sprintf("Error: username and token parameters are both required. [u=%s, p=%s]", *u, *p))
	}

	g := github.NewClient(*u, *p)

	repos, err := g.GetAllRepos("metamx")
	if err != nil {
		panic(err)
	}

	for _, r := range repos {
		fmt.Printf("downloading " + r.Name + "...")
		err = g.DownloadRepoArchive(r, github.AFTarball, "")
		if err != nil {
			fmt.Println("Error: " + err.Error())
		} else {
			fmt.Println("done!")
		}
	}
}
