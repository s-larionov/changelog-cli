package main

import (
	"gitlab.com/3line-io/infrastructure/changelog-cli/pkg/changelog"
)

const (
	BumpManual BumpKind = "manual"
	BumpAuto   BumpKind = "auto"
	BumpPatch  BumpKind = "patch"
	BumpMinor  BumpKind = "minor"
	BumpMajor  BumpKind = "major"
)

var availableKinds = map[BumpKind]struct{}{
	BumpManual: {},
	BumpAuto:   {},
	BumpPatch:  {},
	BumpMinor:  {},
	BumpMajor:  {},
}

var bumpMap = map[changelog.ChangesMajority]BumpKind{
	changelog.NoChanges:    BumpPatch,
	changelog.PatchChanges: BumpPatch,
	changelog.MinorChanges: BumpMinor,
	changelog.MajorChanges: BumpMajor,
}

type BumpKind string
