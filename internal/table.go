package internal

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/google/go-github/v33/github"
)

func printAsTable(repositories []*github.Repository) {
	const padding = 1
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Owner\t Name\t Updated at (UTC)\t Star count")

	for _, repository := range repositories {
		owner := GetOwnerOrOrganisation(repository)
		updated := GetFormattedUpdatedTime(repository)
		name := repository.GetName()
		stars := repository.GetStargazersCount()
		fmt.Fprintln(w, owner, "\t", name, "\t", updated, "\t", stars)
	}

	w.Flush()
}

func GetOwnerOrOrganisation(repository *github.Repository) string {
	var owner string = ""
	repoOrganisation := repository.GetOrganization()
	if repoOrganisation.GetName() != "" {
		owner = repoOrganisation.GetName()
	} else {
		var repoOwner = repository.GetOwner()
		owner = repoOwner.GetLogin()
	}
	return owner
}

func GetFormattedUpdatedTime(repository *github.Repository) string {
	updatedTime := repository.GetUpdatedAt()
	return updatedTime.Time.Format("2006-01-02T15:04:05")
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

func repositoriesToRowValues(repositories []*github.Repository) []rowValues {
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
