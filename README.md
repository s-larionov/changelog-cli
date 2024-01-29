# Command Line Tool for managing CHANGELOG

## How to Build

### Native
```shell
go build -o bin/changelog-cli .
```

### Docker
```shell
docker build -f Dockerfile .\
  -t changelog-cli \
  -e GITHUB_TOKEN=... \
  -e GITLAB_TOKEN=...
```

## How to Use

### Native Commands

Use `--help` flag for reading actual instructions
```shell
./changelog-cli --help
```

#### Show diff between versions

The command prints changes between versions to the STDOUT.

```shell
# Show unreleased changes (the default behaviour)
./changelog-cli [-command=diff] [-from=latest] [-file=CHANGELOG.md]

# Show unreleased changes and non zero exit code on no changes (for using it in the pipelines)
./changelog-cli -fail-on-empty

# Show all changes from the first version to the latest, include unreleased
./changelog-cli [-command=diff] -from=0.0.0 [-file=CHANGELOG.md]

# Show changes between v1.0.0 and v2.0.0
./changelog-cli [-command=diff] -from=1.0.0 -to=2.0.0 [-file=CHANGELOG.md]
```

#### Bump new version:

The command prints updated changelog in Markdown format to STDOUT.

```shell
# Default behaviour:
./changelog-cli -command=bump [-file=CHANGELOG.md]

# Bump to specified version 3.4.9-beta3.12:
./changelog-cli -command=bump -version=3.4.9-beta3.12 [-file=CHANGELOG.md]

# Force bump patch/minor/major version
./changelog-cli -command=bump -bump=minor [-file=CHANGELOG.md]
```

#### Get info about the latest released version:

The command prints latest released version in the changelog to STDOUT.

```shell
# Default behaviour:
./changelog-cli -command=latest_version [-file=CHANGELOG.md]
```

#### Get info about deploy direction:

The command prints deployment direction between versions to STDOUT.

```shell
# Default behaviour:
./changelog-cli -command=direction -from=0.1.2 -to=0.2.0 [-file=CHANGELOG.md]
```

#### Init new changelog:

The command prints default empty changelog to STDOUT.

```shell
# Default behaviour:
./changelog-cli -command=init
```

#### Read file from STDIN:

It's supported in all commands

```shell
# Parse changelog from STDIN:
cat CHANGELOG.md | ./changelog-cli -file=STDIN
```

**Parameters:**
- **command** `string` (default `diff`) \
  Command for execution (`diff`, `bump`, `latest_version`)
- **file** `string` (default `CHANGELOG.md`) \
  Path to the source of the changelog in Markdown format. \
  You can use prefegined value `STDIN` for reading changelog from STDIN.
- **bump** `string` (default `auto`) \
  Specified kind for bumping (`patch`, `minor`, `major`, `auto`)
- **from** `string` (default `latest`) \
  From which version should we generate diff? 
- **to** `string` (default `Unreleased`) \
  Until which version should we generate diff?
- **version** `string` \
  Specified version for bumping. This param will override bump param
- **fail-on-empty** `bool` \
  Pass this parameter if you want trigger an error (non-zero exit code) on no changes on the diff

### Execute Commands inside the Docker
```shell
docker run -v /path/to/CHANGELOG.md:/opt/CHANGELOG.md \
  -e VERSION=11.4.4-alpha1.34 \ 
  -e COMMAND=bump \
  changelog-cli
```

**Parameters**

Parameters will be automatically loaded from the environment variables. Mapping:

- `COMMAND` → `-command=$COMMAND`
- `FILE` → `-file=$FILE`
- `FROM` → `-from=$FROM`
- `TO` → `-to=$TO`
- `BUMP` → `-bump=$BUMP`
- `VERSION` → `-version=$VERSION`
