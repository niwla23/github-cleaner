package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

// fetches repos of authenticated users
func fetchRepos(client github.Client, ctx context.Context) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 20000},
	}
	orgs, _, err := client.Repositories.List(ctx, "", opt)
	return orgs, err
}

// checks if a is a subset of b
func isSubset[E comparable](a []E, b []E) bool {
	for _, itemA := range a {
		if !slices.Contains(b, itemA) {
			return false
		}
	}
	return true
}

func main() {
	ctx := context.Background()
	token, token_set := os.LookupEnv("GH_TOKEN")
	if !token_set {
		panic("No GH_TOKEN found in environment!")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	repos, err := fetchRepos(*client, ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, repo := range repos {
		if repo.GetFork() {
			repoData, _, _ := client.Repositories.Get(ctx, user.GetLogin(), repo.GetName())
			commits, _, _ := client.Repositories.ListCommits(ctx, user.GetLogin(), repo.GetName(), nil)
			parentCommits, _, _ := client.Repositories.ListCommits(ctx, repoData.GetParent().GetOwner().GetLogin(), repoData.GetParent().GetName(), nil)

			commitSHAs := []string{}
			for _, commit := range commits {
				commitSHAs = append(commitSHAs, commit.GetSHA())
			}

			var parentCommitSHAs []string
			for _, commit := range parentCommits {
				parentCommitSHAs = append(parentCommitSHAs, commit.GetSHA())
			}

			if isSubset(commitSHAs, parentCommitSHAs) {
				fmt.Print("repository ", repoData.GetHTMLURL(), " was forked from ", repoData.GetParent().GetFullName())
				fmt.Println(". No Changes were made!")
			}
		}
	}
}
