package pkg

import (
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"

	"gitlab.com/3line-io/infrastructure/changelog-cli/pkg/changelog"
)

const (
	versionLevel     = 2
	changesKindLevel = 3

	prepareContentReplacement = "\n\n$2 "
)

var mdRenderer = md.NewRenderer()
var rePrepareContent = regexp.MustCompile(`(\s*\r?\n)+(#{2,3})\s`)

func ParseMarkdownFile(content []byte) *changelog.Changelog {
	content = rePrepareContent.ReplaceAll(content, []byte(prepareContentReplacement))

	tree := markdown.Parse(content, parser.New())
	header, description, skip := readHeader(tree)
	versions := readVersions(tree, skip)

	cl := changelog.NewChangelog(header, description, versions)

	return cl
}

// readHeader reads all content until first version
func readHeader(tree ast.Node) (header, description string, skip int) {
	for i, node := range tree.GetChildren() {
		if h, ok := node.(*ast.Heading); ok && h.Level == 1 {
			header = renderMarkdownContent(node)
			skip = i + 1
			continue
		}

		if _, ok := node.(*ast.Paragraph); !ok {
			break
		}

		description += renderMarkdownContent(node) + "\n\n"
		skip++
	}

	return strings.TrimSpace(header), strings.TrimSpace(description), skip
}

func readVersions(tree ast.Node, skip int) map[changelog.VersionString]changelog.VersionChanges {
	versions := make(map[changelog.VersionString]changelog.VersionChanges)

	var ver *changelog.Version
	var kind *changelog.ChangesKind
	for _, node := range tree.GetChildren()[skip:] {
		v, ok := isVersion(node)
		if ok {
			ver = &v
			continue
		}

		if ver == nil {
			// For correct Changelog structure it never should happen
			continue
		}

		if _, exist := versions[ver.GetVersion()]; !exist {
			versions[ver.GetVersion()] = changelog.NewVersionChanges(*ver, changelog.NewChanges())
		}

		k, ok := isChangesKind(node)
		if ok {
			kind = &k
			continue
		}

		if kind == nil {
			// For correct Changelog structure it never should happen
			continue
		}

		versions[ver.GetVersion()].Changes.Set(*kind, renderMarkdownContent(node))
	}

	return versions
}

func isVersion(node ast.Node) (changelog.Version, bool) {
	ver, err := changelog.NewVersionFromNode(node, versionLevel)
	if err != nil {
		return changelog.Empty, false
	}

	return ver, true
}

func isChangesKind(node ast.Node) (changelog.ChangesKind, bool) {
	kind, err := changelog.NewChangesKindFromNode(node, changesKindLevel)
	if err != nil {
		return "", false
	}

	return kind, true
}

func renderMarkdownContent(node ast.Node) string {
	return strings.TrimSpace(string(markdown.Render(node, mdRenderer)))
}
