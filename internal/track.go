package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
)

// Track tracks public GitHub repositories, continuously updating according to the given interval.
//
// The given interval must be greater than zero.
func Track(interval time.Duration, token string, minStars int) error {
	for ; ; <-time.Tick(interval) {
		con := context.Background()
		client := github.NewClient(nil)
		if token != "" {
			fmt.Println("Token provided: " + token)
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

		PrintAsTable(filteredRepositories)
	}
}
