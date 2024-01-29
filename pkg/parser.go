package pkg

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	"gitlab.com/3line-io/infrastructure/changelog-cli/pkg/changelog"
)

const (
	versionLevel     = 2
	changesKindLevel = 3

	prepareContentReplacement = "\n\n$2 "
)

var rePrepareContent = regexp.MustCompile(`(\s*\r?\n)+(#{2,3})\s`)

func ParseMarkdownFile(content []byte) *changelog.Changelog {
	content = rePrepareContent.ReplaceAll(content, []byte(prepareContentReplacement))

	tree := goldmark.DefaultParser().Parse(text.NewReader(content))

	header, description, skip := readHeader(content, tree)
	versions := readVersions(content, tree, skip)

	cl := changelog.NewChangelog(header, description, versions)

	return cl
}

// readHeader reads all content until first version
func readHeader(src []byte, tree ast.Node) (header, description string, skip int) {
	i := 0
	for node := tree.FirstChild(); node != nil; node = node.NextSibling() {
		i++

		switch v := node.(type) {
		case *ast.Heading:
			if v.Level == 1 {
				header = renderMarkdownContent(src, node)
				skip = i
				continue
			}
		case *ast.Paragraph:
			description += renderMarkdownContent(src, node) + "\n\n"
			skip++
			continue
			// do nothing
		default:
		}

		break
	}

	return strings.TrimSpace(header), strings.TrimSpace(description), skip
}

func readVersions(src []byte, tree ast.Node, skip int) map[changelog.VersionString]changelog.VersionChanges {
	versions := make(map[changelog.VersionString]changelog.VersionChanges)

	var ver *changelog.Version
	var kind *changelog.ChangesKind

	i := 0
	for node := tree.FirstChild(); node != nil; node = node.NextSibling() {
		i++
		if i <= skip {
			continue
		}

		//for _, node := range tree.GetChildren()[skip:] {
		v, ok := isVersion(src, node)
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

		k, ok := isChangesKind(src, node)
		if ok {
			kind = &k
			continue
		}

		if kind == nil {
			// For correct Changelog structure it never should happen
			continue
		}

		versions[ver.GetVersion()].Changes.Set(*kind, renderMarkdownContent(src, node))
	}

	return versions
}

func isVersion(src []byte, node ast.Node) (changelog.Version, bool) {
	ver, err := changelog.NewVersionFromNode(src, node, versionLevel)
	if err != nil {
		return changelog.Empty, false
	}

	return ver, true
}

func isChangesKind(src []byte, node ast.Node) (changelog.ChangesKind, bool) {
	kind, err := changelog.NewChangesKindFromNode(src, node, changesKindLevel)
	if err != nil {
		return "", false
	}

	return kind, true
}

func renderMarkdownContent(src []byte, node ast.Node) string {
	switch v := node.(type) {
	case *ast.List:
		buf := bytes.NewBufferString("")
		for n := v.FirstChild(); n != nil; n = n.NextSibling() {
			buf.WriteString(fmt.Sprintf("- %s\n", string(n.Text(src))))
		}
		return strings.TrimSpace(buf.String())
	default:
		return string(node.Text(src))
	}
}
