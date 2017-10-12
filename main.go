package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ryanbarry/dodiligence/github"
)

func main() {
	var username = flag.String("username", "", "github username")
	var token = flag.String("token", "", "github personal auth token")
	var org = flag.String("org", "", "github organization to export all repos from")
	var repo = flag.String("repo", "", "name of repository to export (default/unspecified means export all)")
	flag.Parse()

	if len(*username) == 0 || len(*token) == 0 || len(*org) == 0 {
		fmt.Fprintf(os.Stderr, "Error: parameters username, token and org are all required. [username=\"%s\", token=\"%s\", org=\"%s\"]\n", *username, *token, *org)
		os.Exit(1)
	}

	g := github.NewClient(*username, *token)

	var repos []github.GHRepo
	var err error
	if len(*repo) == 0 {
		repos, err = g.GetAllRepos(*org)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when retrieving all repos for \"%s\": %v\n", *org, err)
			os.Exit(2)
		}
	} else {
		r, err := g.GetRepo(*org, *repo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when retrieving \"%s\": %v\n", *org+"/"+*repo, err)
			os.Exit(2)
		}
		repos = append(repos, r)
	}

	for _, r := range repos {
		fmt.Printf("downloading " + r.Name + "...")
		err := g.DownloadRepoArchive(r, github.AFTarball, "")
		if err != nil {
			fmt.Println("Error: " + err.Error())
		} else {
			fmt.Println("done!")
		}
	}
}
