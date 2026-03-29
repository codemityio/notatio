# ![Notatio](logo.jpg)

![coverage-badge-do-not-edit](https://img.shields.io/badge/Coverage-85%25-green.svg?longCache=true&style=flat)

## Table of contents

- [Summary](#summary)
- [Development](#development)
  - [`make`](#make)
- [Installation](#installation)
- [Usage](#usage)
  - [Manual](#manual)
  - [Subcommands](#subcommands)
  - [Docker](#docker)
- [Packages](#packages)
- [License](#license)

## Summary

A tool designed to streamline working with documentation and diagrams.

## Development

To work with the codebase, use `make` command as the primary entry point for all project tools.

Use the arrow keys `↓ ↑ → ←` to navigate the options, and press `/` to toggle search.

### `make`

``` bash
$ make help
build                          Build container image
buildx                         Build container multi platform images and push
check                          Run all CI required targets
cleanup                        Cleanup project
cmd                            Run a command passed as COMMAND= value (e.g. make cmd COMMAND="make check")
cov-open                       Inspect coverage in the browser
cov-report                     Check coverage report
cov                            Check coverage
diff                           Check diff to ensure this project consistency
docs-cmd                       Generate pkg docs
docs-depgraph                  Generate dependency graph
docs-main                      Generate main docs
docs-pkg                       Generate pkg docs
docs-render                    Render diagrams
docs-uml                       Generate UML documentation
docs                           Generate all docs
exec                           Execute built bin (use FLAGS= and COMMAND= environment variables to pass main command flags and subcommand with flags when needed)
fmt                            Format code
gen                            Go generate
go                             Build Go
help                           Prints help for targets with comments
install                        Install binary locally
next                           Create a new version (bump prerelease or patch)
prep                           Prepare dev tools
push                           Push image
reset                          Stop and remove project containers, remove project volumes, remove project images
run-container                  Run container (use FLAGS= and COMMAND= environment variables to pass main command flags and subcommand with flags when needed)
run-go                         Run go (use FLAGS= and COMMAND= environment variables to pass main command flags and subcommand with flags when needed)
statan-fix                     Analyze code and fix
statan                         Analyze code
test-race                      Run race tests
test                           Run tests
update                         Update all dependencies
vendor                         Run go mod vendor
version                        Print the most recent version
```

## Installation

To install the tool use `make install` (directly from the repository clone) or use
`go install github.com/codemityio/notatio@latest`.

> Some of the tools depend on additional commands such as `dot`, `mmdc`, `java` or `chromium-browser`. If any of these
> are missing, you will be notified when using the tools. For the most seamless experience, we recommend using the
> containerized version of this tool.

## Usage

Once you have the tool installed, just use the `notatio` command to get started.

### Manual

``` bash
$ notatio --help
NAME:
   notatio - A new cli application

USAGE:
   notatio [global options] command [command options]

VERSION:
   latest

DESCRIPTION:
   A tool designed to streamline working with documentation and diagrams.

AUTHOR:
   codemityio

COMMANDS:
   coi       
   graphviz  
   mermaid   
   plantuml  
   toc       
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   codemityio
```

### Subcommands

- [`coi`](cmd/coi/README.md) - A simple tool to generate document sections with provided command output.
- [`graphviz`](cmd/graphviz/README.md) - A tool to convert `dot`/`gv` files to `svg`/`png` images.
- [`mermaid`](cmd/mermaid/README.md) - A tool to convert `mmd` files to `svg`/`png` images.
- [`plantuml`](cmd/plantuml/README.md) - A tool to convert `puml` files to `svg`/`png` images.
- [`toc`](cmd/toc/README.md) - A tool to generate table of contents section within a **Markdown** file from a list of
  paths or headers found in a document.

### Docker

``` bash
$ docker run codemityio/notatio
```

## Packages

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
