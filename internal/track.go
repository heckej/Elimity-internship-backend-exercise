package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
)

type columnWidth struct {
	owner, name, updated, stars int
}

type rowValues struct {
	owner, name, updated, stars string
}

func PrintAsTable(repositories []*github.Repository) {
	headerValues := rowValues{
		owner:   "Owner",
		name:    "Name",
		updated: "Updated at (UTC)",
		stars:   "Star count",
	}
	minColWidths := columnWidth{
		owner:   len(headerValues.owner),
		name:    len(headerValues.name),
		updated: len(headerValues.updated),
		stars:   len(headerValues.stars),
	}
	colWidths := calculateColumnWidths(minColWidths, repositories)
	PrintAsTableRow(headerValues, colWidths)
	for _, repository := range repositories {
		var owner string = "unknown"
		repoOrganisation := repository.GetOrganization()
		if repoOrganisation.GetName() != "" {
			owner = repoOrganisation.GetName()
		} else if repoOrganisation.GetLogin() != "" {
			owner = repoOrganisation.GetLogin()
		} else {
			var repoOwner = repository.GetOwner()
			if repoOwner.GetName() != "" {
				owner = repoOwner.GetName()
			} else if repoOwner.GetLogin() != "" {
				owner = repoOwner.GetLogin()
			}
		}
		values := rowValues{
			owner:   owner,
			name:    *repository.Name,
			updated: fmt.Sprint(*repository.UpdatedAt),
			stars:   fmt.Sprint(*repository.StargazersCount),
		}
		PrintAsTableRow(values, colWidths)
	}
}

func PrintAsTableRow(values rowValues, colWidths columnWidth) {
	tableBorder := " | "
	PrintColumn("", values.owner, colWidths.owner, "")
	PrintColumn(tableBorder, values.name, colWidths.name, "")
	PrintColumn(tableBorder, values.updated, colWidths.updated, "")
	PrintColumn(tableBorder, values.stars, colWidths.stars, "\n")
}

func PrintColumn(start string, value string, colWidth int, end string) {
	//formatValueWidth := "%f" + fmt.Sprintf("%d", colWidth) + "v"
	//formatValueWidth := "%f5abc"
	fmt.Printf("%v%-"+fmt.Sprint(colWidth)+"v%v", start, value, end)
}

func calculateColumnWidths(minColWidths columnWidth, repositories []*github.Repository) columnWidth {
	maxOwnerLength := minColWidths.owner
	maxNameLength := minColWidths.name
	maxUpdatedLength := minColWidths.updated
	maxStarsLength := minColWidths.stars
	for _, repository := range repositories {
		var owner string = "unknown"
		repoOrganisation := repository.GetOrganization()
		if repoOrganisation.GetName() != "" {
			owner = repoOrganisation.GetName()
		} else if repoOrganisation.GetLogin() != "" {
			owner = repoOrganisation.GetLogin()
		} else {
			var repoOwner = repository.GetOwner()
			if repoOwner.GetName() != "" {
				owner = repoOwner.GetName()
			} else if repoOwner.GetLogin() != "" {
				owner = repoOwner.GetLogin()
			}
		}
		maxOwnerLength = max(maxOwnerLength, len(owner))
		maxNameLength = max(maxNameLength, len(*repository.Name))
		maxUpdatedLength = max(maxUpdatedLength, len(fmt.Sprint(*repository.UpdatedAt)))
		maxStarsLength = max(maxStarsLength, len(fmt.Sprint(*repository.StargazersCount)))
	}
	colWidths := columnWidth{
		owner:   maxOwnerLength,
		name:    maxNameLength,
		updated: maxUpdatedLength,
		stars:   maxStarsLength,
	}
	return colWidths
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

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
