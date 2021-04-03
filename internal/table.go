package internal

import (
	"fmt"

	"github.com/google/go-github/v33/github"
)

func printAsTable(repositories []*github.Repository) {
	headerValues := rowValues{
		owner:   "Owner",
		name:    "Name",
		updated: "Updated at (UTC)",
		stars:   "Star count",
	}
	rows := []rowValues{headerValues}
	repositoryValues := repositoriesToValues(repositories)
	rows = append(rows, repositoryValues...)
	colWidths := calculateColumnWidths(rows)
	for _, row := range rows {
		printAsTableRow(row, colWidths)
	}
}

func printAsTableRow(values rowValues, colWidths columnWidth) {
	tableBorder := " | "
	printColumn("", values.owner, colWidths.owner, "")
	printColumn(tableBorder, values.name, colWidths.name, "")
	printColumn(tableBorder, values.updated, colWidths.updated, "")
	printColumn(tableBorder, values.stars, colWidths.stars, "\n")
}

func printColumn(start string, value string, colWidth int, end string) {
	fmt.Printf("%v%-"+fmt.Sprint(colWidth)+"v%v", start, value, end)
}

func repositoriesToValues(repositories []*github.Repository) []rowValues {
	repositoryValues := []rowValues{}
	for _, repository := range repositories {
		var owner string = ""
		repoOrganisation := repository.GetOrganization()
		if repoOrganisation.GetName() != "" {
			owner = repoOrganisation.GetName()
		} else {
			var repoOwner = repository.GetOwner()
			owner = repoOwner.GetLogin()
		}
		updatedTime := repository.GetUpdatedAt()
		updated := updatedTime.Time.Format("2006-01-02T15:04:05")
		values := rowValues{
			owner:   owner,
			name:    repository.GetName(),
			updated: updated,
			stars:   fmt.Sprint(repository.GetStargazersCount()),
		}
		repositoryValues = append(repositoryValues, values)
	}
	return repositoryValues
}

func calculateColumnWidths(rows []rowValues) columnWidth {
	maxOwnerLength, maxNameLength, maxUpdatedLength, maxStarsLength := 0, 0, 0, 0
	for _, row := range rows {
		maxOwnerLength = max(maxOwnerLength, len(row.owner))
		maxNameLength = max(maxNameLength, len(row.name))
		maxUpdatedLength = max(maxUpdatedLength, len(row.updated))
		maxStarsLength = max(maxStarsLength, len(row.stars))
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

type columnWidth struct {
	owner, name, updated, stars int
}

type rowValues struct {
	owner, name, updated, stars string
}
