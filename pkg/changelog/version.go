package changelog

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/yuin/goldmark/ast"
)

const UnreleasedValue VersionString = "Unreleased"
const LatestValue VersionString = "Latest"
const day = 24 * time.Hour

var (
	Empty      = Version{}
	Unreleased = RequireVersionFromString(UnreleasedValue, nil)
	Latest     = RequireVersionFromString(LatestValue, nil)
)

var (
	ErrNotIsVersion = errors.New("the node is not version")
)

var re = regexp.MustCompile(`^(\[(.+?)]|(.+?))(\s-\s(\d{4}-\d{2}-\d{2}))?$`)

type VersionString string

type Version struct {
	version       VersionString
	parsedVersion *semver.Version
	date          time.Time
}

func NewVersion(version VersionString, date *time.Time) (Version, error) {
	ver := Version{
		version: version,
	}

	if ver.IsCommon() {
		parsed, err := semver.NewVersion(string(ver.version))
		if err != nil {
			return Empty, fmt.Errorf("%v: %v", ErrNotIsVersion, err)
		}
		ver.version = VersionString(parsed.String())
		ver.parsedVersion = parsed
	}

	if date != nil {
		*date = date.Truncate(day)
		ver.date = *date
	}

	return ver, nil
}

func RequireVersionFromString(version VersionString, date *time.Time) Version {
	ver, _ := NewVersion(version, date)

	return ver
}

// NewVersionFromNode is method for parsing the version from changelog in markdown format
//
// Valid variations:
// - [version] - 2000-01-01
// - [version]
// - version - 2000-01-01
// - version
//
// The version should be supported by semver or be constant "Unrealized"
func NewVersionFromNode(src []byte, node ast.Node, requiredLevel int) (Version, error) {
	h, ok := node.(*ast.Heading)
	if !ok {
		return Empty, ErrNotIsVersion
	}

	if h.Level != requiredLevel {
		return Empty, ErrNotIsVersion
	}

	text := strings.TrimSpace(string(h.Text(src)))

	matches := re.FindAllStringSubmatch(text, 1)
	if matches == nil {
		return Empty, ErrNotIsVersion
	}

	ver := matches[0][2]
	if ver == "" {
		ver = matches[0][3]
	}

	date, err := time.Parse("2006-01-02", matches[0][5])
	if err == nil {
		return NewVersion(VersionString(ver), &date)
	}

	return NewVersion(VersionString(ver), nil)
}

func (v Version) IsValid() bool {
	return v.version != ""
}

func (v Version) GetVersion() VersionString {
	return v.version
}

func (v Version) IsUnrealized() bool {
	return strings.EqualFold(string(v.version), string(UnreleasedValue))
}

func (v Version) IsLatest() bool {
	return strings.EqualFold(string(v.version), string(LatestValue))
}

func (v Version) IsCommon() bool {
	return !v.IsUnrealized() && !v.IsLatest()
}

func (v Version) GetDate() time.Time {
	return v.date
}

func (v Version) LessThen(ver Version) bool {
	isVUnrealized := v.IsUnrealized()
	isVerUnrealized := ver.IsUnrealized()

	if isVUnrealized && !isVerUnrealized {
		return false
	} else if isVerUnrealized && !isVUnrealized {
		return true
	}

	isVLatest := v.IsLatest()
	isVerLatest := ver.IsLatest()

	if isVLatest && !isVerLatest {
		return false
	} else if isVerLatest && !isVLatest {
		return true
	}

	return v.parsedVersion.LessThan(ver.parsedVersion)
}

func (v Version) GreaterThan(ver Version) bool {
	isVUnrealized := v.IsUnrealized()
	isVerUnrealized := ver.IsUnrealized()

	if isVUnrealized && !isVerUnrealized {
		return true
	} else if isVerUnrealized {
		return false
	}

	isVLatest := v.IsLatest()
	isVerLatest := ver.IsLatest()

	if isVLatest && !isVerLatest {
		return true
	} else if isVerLatest {
		return false
	}

	return v.parsedVersion.GreaterThan(ver.parsedVersion)
}

func (v Version) Equal(ver Version) bool {
	isVUnrealized := v.IsUnrealized()
	isVerUnrealized := ver.IsUnrealized()

	if isVUnrealized && isVerUnrealized {
		return true
	} else if isVerUnrealized || isVUnrealized {
		return false
	}

	isVLatest := v.IsLatest()
	isVerLatest := ver.IsLatest()

	if isVLatest && isVerLatest {
		return true
	} else if isVerLatest || isVLatest {
		return false
	}

	return v.parsedVersion.Equal(ver.parsedVersion)
}

func (v Version) BumpMajor() Version {
	if !v.IsCommon() {
		return v
	}

	if !v.IsValid() {
		return Empty
	}

	now := time.Now()
	bumped := v.parsedVersion.IncMajor()

	return RequireVersionFromString(VersionString(bumped.String()), &now)
}

func (v Version) BumpMinor() Version {
	if !v.IsCommon() {
		return v
	}

	if !v.IsValid() {
		return Empty
	}

	now := time.Now()
	bumped := v.parsedVersion.IncMinor()

	return RequireVersionFromString(VersionString(bumped.String()), &now)
}

func (v Version) BumpPatch() Version {
	if !v.IsCommon() {
		return v
	}

	if !v.IsValid() {
		return Empty
	}

	now := time.Now()
	bumped := v.parsedVersion.IncPatch()

	return RequireVersionFromString(VersionString(bumped.String()), &now)
}
