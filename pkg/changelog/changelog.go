package changelog

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	ErrVersionAlreadyExists = errors.New("version is already exist")
	ErrNothingToRelease     = errors.New("changelog does not contain unreleased changes")
)

type Changelog struct {
	Header      string
	Description string
	Versions    map[VersionString]VersionChanges
}

func NewChangelog(header, description string, versions map[VersionString]VersionChanges) *Changelog {
	return &Changelog{
		Header:      header,
		Description: description,
		Versions:    versions,
	}
}

func (l *Changelog) GetChanges(ver Version) (Changes, bool) {
	changes, ok := l.Versions[ver.GetVersion()]

	return changes.Changes, ok
}

func (l *Changelog) GetLatestVersion() Version {
	ver := RequireVersionFromString("0.0", nil)

	for _, changes := range l.Versions {
		if changes.Version.IsUnrealized() {
			continue
		}

		if changes.Version.GreaterThan(ver) {
			ver = changes.Version
		}
	}

	return ver
}

func (l *Changelog) Release(ver Version) error {
	if _, ok := l.Versions[ver.GetVersion()]; ok {
		return fmt.Errorf("%v: %s", ErrVersionAlreadyExists, ver.GetVersion())
	}

	changes, ok := l.GetChanges(Unreleased)
	if !ok || changes.GetMajority() == NoChanges {
		return ErrNothingToRelease
	}

	l.Versions[Unreleased.GetVersion()] = NewVersionChanges(Unreleased, NewChanges())
	l.Versions[ver.GetVersion()] = NewVersionChanges(ver, changes)

	return nil
}

func (l *Changelog) Add(ver Version, changes Changes) error {
	if _, ok := l.Versions[ver.GetVersion()]; ok {
		return fmt.Errorf("%v: %s", ErrVersionAlreadyExists, ver.GetVersion())
	}

	l.Versions[ver.GetVersion()] = NewVersionChanges(ver, changes)

	return nil
}

func (l *Changelog) GetSortedVersions() []Version {
	versions := make([]Version, 0, len(l.Versions))
	for _, changes := range l.Versions {
		versions = append(versions, changes.Version)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].GreaterThan(versions[j])
	})

	return versions
}

func (l *Changelog) GetDiff(from, to Version) Changes {
	diff := NewChanges()

	if from.GreaterThan(to) {
		from, to = to, from
	}

	latest := l.GetLatestVersion()
	if from.IsLatest() {
		from = latest
	}
	if to.IsLatest() {
		to = latest
	}

	for _, ver := range l.GetSortedVersions() {
		if ver.IsLatest() {
			ver = latest
		}

		if !ver.GreaterThan(from) {
			continue
		}

		if ver.GreaterThan(to) {
			continue
		}

		changes, _ := l.GetChanges(ver)
		for kind, details := range changes {
			union := diff.Get(kind)
			union += "\n"
			union += details
			union = strings.TrimSpace(union)

			diff.Set(kind, union)
		}
	}

	return diff
}

func (l *Changelog) ToMarkdown() string {
	output := l.Header + "\n\n" + l.Description + "\n\n"

	for _, ver := range l.GetSortedVersions() {
		versionString := fmt.Sprintf("## [%s] - %s", ver.GetVersion(), ver.GetDate().Format("2006-01-02"))
		if ver.IsUnrealized() {
			versionString = fmt.Sprintf("## [%s]", ver.GetVersion())
		}

		output += versionString + "\n\n"

		changes, _ := l.GetChanges(ver)
		output += changes.ToMarkdown() + "\n\n"
	}

	return strings.TrimSpace(output)
}
