package main

import (
	"fmt"
	"os"

	"github.com/s-larionov/changelog-cli/pkg/changelog"
)

const (
	Upgrade  = "UPGRADE"
	Rollback = "ROLLBACK"
	Redeploy = "REDEPLOY"
)

func getDirectionCommand(cl *changelog.Changelog) {
	if to.IsUnrealized() || from.IsUnrealized() {
		Usage("You have to specified 'from' and 'to' versions instead of using UNRELEASED keyword")
		os.Exit(1)
	}

	// If keyword "latest" is used replace it to real latest version from Changelog
	latest := cl.GetLatestVersion()
	if from.IsLatest() {
		from = latest
	}
	if to.IsLatest() {
		to = latest
	}

	if _, fromExists := cl.GetChanges(from); !fromExists {
		_, _ = fmt.Fprintf(os.Stderr, "[WARN] Version %s does not exist in CHANGELOG.md\n", from.GetVersion())
	}

	if _, toExists := cl.GetChanges(to); !toExists {
		_, _ = fmt.Fprintf(os.Stderr, "[WARN] Version %s does not exist in CHANGELOG.md\n", to.GetVersion())
	}

	var direction string
	switch {
	case to.GreaterThan(from):
		direction = Upgrade
	case to.LessThen(from):
		direction = Rollback
	default:
		direction = Redeploy
	}

	fmt.Println(direction)
	os.Exit(0)
}
