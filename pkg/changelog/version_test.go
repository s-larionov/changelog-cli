package changelog

import (
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/smartystreets/goconvey/convey"
)

func TestNewFromNode(t *testing.T) {
	versions := [][]string{
		{"## [1.0.2] - 2022-01-04", "1.0.2", "2022-01-04"},
		{"## [1.2.3-beta.1+build345]", "1.2.3-beta.1+build345", "0001-01-01"},
		{"## 0.4 - 2002-02-07", "0.4.0", "2002-02-07"},
		{"## 1.0.2-patch2", "1.0.2-patch2", "0001-01-01"},
		{"## Unreleased", "Unreleased", "0001-01-01"},
		{"## Latest", "Latest", "0001-01-01"},
	}

	for _, v := range versions {
		node := markdown.Parse([]byte(v[0]), parser.New()).GetChildren()[0]

		convey.Convey(v[0], t, func() {
			ver, err := NewVersionFromNode(node, 2)

			convey.So(err, convey.ShouldBeNil)
			convey.So(ver.IsValid(), convey.ShouldBeTrue)
			convey.So(ver.GetVersion(), convey.ShouldNotBeEmpty)
			convey.So(ver.GetVersion(), convey.ShouldEqual, v[1])
			convey.So(ver.GetDate().Format("2006-01-02"), convey.ShouldEqual, v[2])
		})
	}
}
