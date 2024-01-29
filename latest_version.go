package main

import (
	"fmt"

	"gitlab.com/3line-io/infrastructure/changelog-cli/pkg/changelog"
)

func latestVersionCommand(cl *changelog.Changelog) {
	fmt.Println(cl.GetLatestVersion().GetVersion())
}
