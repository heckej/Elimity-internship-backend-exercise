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
		fmt.Fprintln(w,
			getOwnerOrOrganisation(repository), "\t",
			repository.GetName(), "\t",
			getFormattedUpdatedTime(repository), "\t",
			repository.GetStargazersCount())
	}

	w.Flush()
}

func getOwnerOrOrganisation(repository *github.Repository) string {
	repoOrganisation := repository.GetOrganization()
	if repoOrganisation.GetName() != "" {
		return repoOrganisation.GetName()
	} else {
		repoOwner := repository.GetOwner()
		return repoOwner.GetLogin()
	}
}

func getFormattedUpdatedTime(repository *github.Repository) string {
	updatedTime := repository.GetUpdatedAt()
	return updatedTime.Time.Format("2006-01-02T15:04:05")
}
