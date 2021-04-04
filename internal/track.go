package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

// Track tracks public GitHub repositories that have at least minStars stars, continuously updating according to the given interval.
//
// The given interval must be greater than zero.
// If the given token is non-empty, it is used for authenticated requests. Otherwise, it is ignored.
// The given minStars must be non-negative.
func Track(interval time.Duration, token string, minStars int) error {
	for ; ; <-time.Tick(interval) {
		con := context.Background()
		client := github.NewClient(nil)
		if token != "" {
			fmt.Println("Token provided. Using authenticated requests.")
			tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
			authenticationClient := oauth2.NewClient(con, tokenSource)
			client = github.NewClient(authenticationClient)
		} else {
			fmt.Println("No token provided. Using anonymous requests.")
		}
		fmt.Println("Selecting updated repo's with >=", minStars, "stars.")
		listOptions := github.ListOptions{PerPage: 3}
		searchOptions := &github.SearchOptions{ListOptions: listOptions, Sort: "updated"}
		result, _, err := client.Search.Repositories(con, "is:public", searchOptions)
		if err != nil {
			return err
		}

		filteredRepositories := []*github.Repository{}

		for _, repository := range result.Repositories {
			if *repository.StargazersCount >= minStars {
				filteredRepositories = append(filteredRepositories, repository)
			}
		}

		printAsTable(filteredRepositories)
	}
}
