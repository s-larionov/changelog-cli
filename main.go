package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"gitlab.com/3line-io/infrastructure/changelog-cli/pkg"
	"gitlab.com/3line-io/infrastructure/changelog-cli/pkg/changelog"
)

const (
	InitCommand          Command = "init"
	DiffCommand          Command = "diff"
	BumpCommand          Command = "bump"
	LatestVersionCommand Command = "latest_version"
	GetDirectionCommand  Command = "direction"

	UseSTDIN = "stdin"
)

type Command string

var (
	command              Command
	filepath             string
	manualVersion        changelog.Version
	bump                 BumpKind
	fromString, toString string
	from, to             changelog.Version
	failOnEmpty          bool
)

func init() {
	flag.Usage = func() {
		Usage("")
	}

	flag.StringVar(&filepath, "file", "CHANGELOG.md", "Path to the source of the changelog in markdown format or 'STDIN' for reading content from STDIN")
	flag.StringVar(&fromString, "from", "latest", "From which version should we generate diff?")
	flag.StringVar(&toString, "to", "Unreleased", "Until which version should we generate diff?")
	flag.BoolVar(&failOnEmpty, "fail-on-empty", false, "If this param is passed the tool will return non-zero exit code on 'no changes'")
	commandStr := flag.String("command", "diff", "Command for execution (diff, bump, latest_version, direction, init)")
	bumpSrc := flag.String("bump", "auto", "Specified kind for bumping (patch, minor, major, auto)")
	versionSrc := flag.String("version", "", "Specified version for bumping. This param will override bump param")

	flag.Parse()

	command = Command(strings.ToLower(*commandStr))
	if command == InitCommand {
		return
	}

	if filepath == "" {
		Usage("Empty filename is passed")
		os.Exit(1)
	}

	switch command {
	case DiffCommand, GetDirectionCommand:
		var err error
		from, err = changelog.NewVersion(changelog.VersionString(fromString), nil)
		if err != nil {
			Usage(fmt.Sprintf("Wrong format for 'from' version: %v\n", err))
			os.Exit(1)
		}

		to, err = changelog.NewVersion(changelog.VersionString(toString), nil)
		if err != nil {
			Usage(fmt.Sprintf("Wrong format for 'to' version: %v\n", err))
			os.Exit(1)
		}
	case BumpCommand:
		if _, ok := availableKinds[BumpKind(strings.ToLower(*bumpSrc))]; !ok {
			Usage(fmt.Sprintf("Wrong bump parameter: %v\n", *bumpSrc))
			os.Exit(1)
		}
		bump = BumpKind(*bumpSrc)

		if *versionSrc != "" {
			var err error
			manualVersion, err = changelog.NewVersion(changelog.VersionString(*versionSrc), nil)
			if err != nil {
				Usage(fmt.Sprintf("Wrong format for to-version: %v\n", err))
				os.Exit(1)
			}

			bump = BumpManual
		}
	case LatestVersionCommand:
	default:
		Usage(fmt.Sprintf("Wrong command: %v\n", *commandStr))
		os.Exit(1)
	}
}

func readChangelog(filepath string) ([]byte, error) {
	if !strings.EqualFold(filepath, UseSTDIN) {
		return os.ReadFile(filepath)
	}

	return io.ReadAll(os.Stdin)
}

func main() {
	if command == InitCommand {
		initCommand()
		return
	}

	clContent, err := readChangelog(filepath)
	if err != nil {
		Usage(fmt.Sprintf("Unable to read changelog file: %v\n", err))
		os.Exit(1)
	}

	cl := pkg.ParseMarkdownFile(clContent)

	switch command {
	case DiffCommand:
		diffCommand(cl)
	case BumpCommand:
		bumpCommand(cl)
	case LatestVersionCommand:
		latestVersionCommand(cl)
	case GetDirectionCommand:
		getDirectionCommand(cl)
	}
}

func Usage(msg string) {
	if msg != "" {
		fmt.Println(msg)
		fmt.Println()
	}

	fmt.Printf("Usage of %s:\n", os.Args[0])
	fmt.Println("  Show diff between versions:")
	fmt.Printf("    %s -command=diff [-file=CHANGELOG.md] [-from=latest] [-to=Unreleased]\n", os.Args[0])
	fmt.Println()
	fmt.Println("  Bump new version:")
	fmt.Printf("    %s -command=bump [-file=CHANGELOG.md] [-bump=auto] [-version=]\n", os.Args[0])
	fmt.Println()
	fmt.Println("  Init new changelog:")
	fmt.Printf("    %s -command=init\n", os.Args[0])
	fmt.Println()
	fmt.Println("  Show release direction (UPGRADE, ROLLBACK, REDEPLOY):")
	fmt.Printf("    %s -command=direction [-file=CHANGELOG.md] -from=0.1.4 -to=2.3.4\n", os.Args[0])
	fmt.Println()
	fmt.Println("  Get the latest released version from the CHANGELOG:")
	fmt.Printf("    %s -command=latest_version [-file=CHANGELOG.md]\n", os.Args[0])
	fmt.Println()

	fmt.Println("Parameters:")
	flag.PrintDefaults()
}
