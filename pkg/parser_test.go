package pkg

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/s-larionov/changelog-cli/pkg/changelog"
)

func TestParseMarkdownFile(t *testing.T) {
	const content = "## Unreleased\n### Added\n- line 1\n- line 2\n### Fixed\n- line 1\n- line 2\n## 0.1.0\n### Added\n- line 1\n- line 2\n### Changed\n- line 1\n- line 2"

	convey.Convey("parsing changelog", t, func() {
		cl := ParseMarkdownFile([]byte(content))

		convey.Convey("should be successful", func() {
			convey.So(cl.Header, convey.ShouldBeEmpty)
			convey.So(cl.Description, convey.ShouldBeEmpty)

			convey.Convey("and should contains correct versions", func() {
				convey.So(cl.Versions, convey.ShouldHaveLength, 2)
				convey.So(cl.GetLatestVersion().GetVersion(), convey.ShouldEqual, "0.1.0")

				unreleased, ok := cl.GetChanges(changelog.Unreleased)
				convey.So(ok, convey.ShouldBeTrue)
				convey.So(unreleased.ToMarkdown(), convey.ShouldEqual, "### Fixed\n- line 1\n- line 2\n\n### Added\n- line 1\n- line 2")
			})
		})
	})
}

func TestParseMarkdownFile_ToMarkdown(t *testing.T) {
	const md = "## Unreleased\n### Added\n- changes 1\n- changes 2"

	convey.Convey("parsing changelog", t, func() {
		cl := ParseMarkdownFile([]byte(md))

		convey.Convey("should be successful", func() {
			convey.So(cl.Header, convey.ShouldBeEmpty)
			convey.So(cl.Description, convey.ShouldBeEmpty)

			convey.Convey("and rendered version should be in braces []", func() {
				convey.So(cl.ToMarkdown(), convey.ShouldEqual, "## [Unreleased]\n\n### Added\n- changes 1\n- changes 2")
			})
		})
	})
}
