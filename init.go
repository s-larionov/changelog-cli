package main

import (
	"fmt"

	"github.com/s-larionov/changelog-cli/pkg/changelog"
)

const (
	clDefaultHeader = `# Changelog

All notable changes to this project will be documented in this file.`

	clDefaultDescription = `The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).`

	clDefaultAddChangelogChanges = "- Add CHANGELOG.md"
)

func initCommand() {
	versions := make(map[changelog.VersionString]changelog.VersionChanges)
	cl := changelog.NewChangelog(clDefaultHeader, clDefaultDescription, versions)
	changes := changelog.NewChanges()
	changes.Set(changelog.Added, clDefaultAddChangelogChanges)
	_ = cl.Add(changelog.Unreleased, changes)

	fmt.Println(cl.ToMarkdown())
}
