# ![Notatio](logo.jpg)

![coverage-badge-do-not-edit](https://img.shields.io/badge/Coverage-5%25-red.svg?longCache=true&style=flat)

## Table of contents

- [Summary](#summary)
- [Development](#development)
- [Installation](#installation)
- [Usage](#usage)
  - [Manual](#manual)
  - [Subcommands](#subcommands)
- [Packages](#packages)
- [License](#license)

## Summary

A tool designed to streamline working with documentation and diagrams.

## Development

To work with the codebase, use `make` command as the primary entry point for all project tools.

Use the arrow keys `↓ ↑ → ←` to navigate the options, and press `/` to toggle search.

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
   v0.0.0

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

- [`plantuml`](cmd/plantuml/README.md) - A tool to convert `puml` files to `svg`/`png` images.
- [`toc`](cmd/toc/README.md) - A tool to generate table of contents section within a **Markdown** file from a list of
  paths or headers found in a document.
- [`coi`](cmd/coi/README.md) - A simple tool to generate document sections with provided command output.
- [`graphviz`](cmd/graphviz/README.md) - A tool to convert `dot`/`gv` files to `svg`/`png` images.
- [`mermaid`](cmd/mermaid/README.md) - A tool to convert `mmd` files to `svg`/`png` images.

## Packages

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
