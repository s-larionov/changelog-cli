package main

import (
	"fmt"
	"os"

	"gitlab.com/3line-io/infrastructure/changelog-cli/pkg/changelog"
)

func diffCommand(cl *changelog.Changelog) {
	if from.IsLatest() {
		from = cl.GetLatestVersion()
	}

	changes := cl.GetDiff(from, to)

	// If from and to versions are the same then diff between them is changes in exactly this version
	if from.Equal(to) {
		changes, _ = cl.GetChanges(to)
	}

	output := changes.ToMarkdown()

	if output != "" {
		fmt.Println(output)
	}

	if output == "" && failOnEmpty {
		os.Exit(1)
	}
}
