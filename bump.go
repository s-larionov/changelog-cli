package main

import (
	"fmt"
	"os"

	"gitlab.com/3line-io/infrastructure/changelog-cli/pkg/changelog"
)

func bumpCommand(cl *changelog.Changelog) {
	unreleased, ok := cl.GetChanges(changelog.Unreleased)
	if !ok {
		Usage("Changelog does not contain unreleased changes")
		os.Exit(1)
	}

	if bump == BumpAuto {
		majority := unreleased.GetMajority()
		if majority == changelog.NoChanges {
			Usage("Changelog does not contain unreleased changes")
			os.Exit(1)
		}

		bump = bumpMap[majority]
	}

	latestVersion := cl.GetLatestVersion()
	var version changelog.Version
	switch bump {
	case BumpManual:
		version = manualVersion
	case BumpPatch:
		version = latestVersion.BumpPatch()
	case BumpMinor:
		version = latestVersion.BumpMinor()
	case BumpMajor:
		version = latestVersion.BumpMajor()
	}

	if err := cl.Release(version); err != nil {
		Usage(fmt.Sprintf("Unable to make release: %v", err))
		os.Exit(1)
	}

	fmt.Println(cl.ToMarkdown())
}
