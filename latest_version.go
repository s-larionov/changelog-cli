package main

import (
	"fmt"

	"github.com/s-larionov/changelog-cli/pkg/changelog"
)

func latestVersionCommand(cl *changelog.Changelog) {
	fmt.Println(cl.GetLatestVersion().GetVersion())
}
