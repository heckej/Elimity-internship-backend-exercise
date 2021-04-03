package internal

import (
	"fmt"

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
	rows := []rowValues{headerValues}
	repositoryValues := RepositoriesToValues(repositories)
	rows = append(rows, repositoryValues...)
	colWidths := CalculateColumnWidths(rows)
	for _, row := range rows {
		PrintAsTableRow(row, colWidths)
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
	fmt.Printf("%v%-"+fmt.Sprint(colWidth)+"v%v", start, value, end)
}

func RepositoriesToValues(repositories []*github.Repository) []rowValues {
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

func CalculateColumnWidths(rows []rowValues) columnWidth {
	maxOwnerLength, maxNameLength, maxUpdatedLength, maxStarsLength := 0, 0, 0, 0
	for _, row := range rows {
		maxOwnerLength = Max(maxOwnerLength, len(row.owner))
		maxNameLength = Max(maxNameLength, len(row.name))
		maxUpdatedLength = Max(maxUpdatedLength, len(row.updated))
		maxStarsLength = Max(maxStarsLength, len(row.stars))
	}
	colWidths := columnWidth{
		owner:   maxOwnerLength,
		name:    maxNameLength,
		updated: maxUpdatedLength,
		stars:   maxStarsLength,
	}
	return colWidths
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
