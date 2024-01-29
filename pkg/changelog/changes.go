package changelog

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown/ast"
)

const (
	Added      ChangesKind = "Added"      // for new features.
	Changed    ChangesKind = "Changed"    // for changes in existing functionality.
	Deprecated ChangesKind = "Deprecated" // for soon-to-be removed features.
	Removed    ChangesKind = "Removed"    // for now removed features.
	Fixed      ChangesKind = "Fixed"      // for any bug fixes.
	Security   ChangesKind = "Security"   // in case of vulnerabilities.
)

const (
	NoChanges ChangesMajority = iota
	PatchChanges
	MinorChanges
	MajorChanges
)

var OrderedKinds = []ChangesKind{
	Security,
	Fixed,
	Added,
	Changed,
	Removed,
	Deprecated,
}

var MajorityMap = map[ChangesKind]ChangesMajority{
	Security:   PatchChanges,
	Fixed:      PatchChanges,
	Deprecated: PatchChanges,
	Added:      MinorChanges,
	Changed:    MinorChanges,
	Removed:    MajorChanges,
}

var ErrNotIsChangesKind = errors.New("the node is not kind of changes")

type ChangesKind string

type ChangesMajority uint

func NewChangesKindFromNode(node ast.Node, requiredLevel int) (ChangesKind, error) {
	h, ok := node.(*ast.Heading)
	if !ok {
		return "", ErrNotIsChangesKind
	}

	if h.Level != requiredLevel {
		return "", ErrNotIsChangesKind
	}

	if len(h.GetChildren()) != 1 {
		return "", ErrNotIsChangesKind
	}

	text := strings.TrimSpace(string(h.GetChildren()[0].AsLeaf().Literal))

	if text == "" {
		return "", ErrNotIsChangesKind
	}

	return ChangesKind(text), nil
}

type VersionChanges struct {
	Version Version
	Changes Changes
}

func NewVersionChanges(ver Version, changes Changes) VersionChanges {
	return VersionChanges{
		Version: ver,
		Changes: changes,
	}
}

type Changes map[ChangesKind]string

func NewChanges() Changes {
	return make(Changes)
}

func (c Changes) Set(kind ChangesKind, changes string) {
	if c == nil {
		return
	}

	c[kind] = changes
}

func (c Changes) Get(kind ChangesKind) string {
	if c == nil {
		return ""
	}

	return c[kind]
}

func (c Changes) Has(kind ChangesKind) bool {
	if c == nil {
		return false
	}

	changes, ok := c[kind]
	if !ok {
		return false
	}

	if changes == "" {
		return false
	}

	return true
}

func (c Changes) GetMajority() ChangesMajority {
	if len(c) == 0 {
		return NoChanges
	}

	majority := NoChanges
	for kind, m := range MajorityMap {
		if c.Has(kind) && m > majority {
			majority = m
		}
	}

	return majority
}

func (c Changes) ToMarkdown() string {
	output := ""

	for _, kind := range OrderedKinds {
		if !c.Has(kind) {
			continue
		}

		output += fmt.Sprintf("### %s\n", kind)
		output += fmt.Sprintf("%s\n\n", c.Get(kind))
	}

	return strings.TrimSpace(output)
}
