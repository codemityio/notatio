# `toc`

## Table of contents

- [Summary](#summary)
- [Manual](#manual)
- [Subcommands](#subcommands)
  - [`int`](#int)
  - [`ext`](#ext)
- [Usage](#usage)
  - [`int`](#int)
    - [Original Markdown file content](#original-markdown-file-content)
    - [Command](#command)
    - [Result Markdown file content](#result-markdown-file-content)

## Summary

A tool to generate table of contents section within a **Markdown** file from a list of paths or headers found in a
document.

## Manual

``` bash
$ notatio toc --help
NAME:
   notatio toc

USAGE:
   notatio toc [command options]

DESCRIPTION:
   Table of contents generator.

COMMANDS:
   int      Generate table of content from headers within a document 
              e.g. toc --document-path=README.md --header="Table of contents" --limiter-right="##" int.
   ext      Generate table of content within a document and use provided paths as a list 
              e.g. toc --document-path=README.md --header="Table of contents" --limiter-right="##" ext --path=one/document.md --path=two/document.md.
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --document-path value  markdown file path to be updated
   --header value         header to use for document lookups and generation
   --limiter-left value   string to use as a lookup limiter (default: "##")
   --limiter-right value  string to use as a lookup limiter - empty will use end of file as a limit (default: "##")
   --index value          index of a section to be used as a placeholder (useful if limiters refer to more than one section,
      0 = replace all) (default: 0)
   --help, -h  show help
```

## Subcommands

### `int`

``` bash
$ notatio toc --document-path=README.md --header='Table of contents' int --help
NAME:
   notatio toc int - Generate table of content from headers within a document 
                       e.g. toc --document-path=README.md --header="Table of contents" --limiter-right="##" int.

USAGE:
   notatio toc int [command options]

OPTIONS:
   --start-from-level value  indicate what level of headers to start from (default: 0)
   --start-from-item value   indicate what item from the list to start from (default: 0)
   --help, -h                show help
```

### `ext`

``` bash
$ notatio toc --document-path=README.md --header='Table of contents' ext --help
NAME:
   notatio toc ext - Generate table of content within a document and use provided paths as a list 
                       e.g. toc --document-path=README.md --header="Table of contents" --limiter-right="##" ext --path=one/document.md --path=two/document.md.

USAGE:
   notatio toc ext [command options]

OPTIONS:
   --path value [ --path value ]  path to be included in the table of contents
   --summary-header value         summary header to use for document lookups
   --summary-limiter-left value   string to use as a summary lookup limiter (default: "##")
   --summary-limiter-right value  string to use as a summary lookup limiter - empty will use end of file as a limit (default: "##")
   --help, -h                     show help
```

## Usage

### `int`

#### Original Markdown file content

``` markdown
...

## Table of contents

## Summary

...
```

#### Command

The following command will take all provided paths and generate a list of links within a document.

``` shell
notatio toc --document-path=README.md --header="Table of contents" --limiter-right=## ext --path=one/one.md --path=two/two.md
```

#### Result Markdown file content

``` markdown
...

## Table of contents

  - [Summary](#summary)
  - [Manual](#manual)
  - [Usage](#usage)
  - [Example](#example)
    - [Original Markdown file content](#original-markdown-file-content)
    - [Command](#command)
    - [Result Markdown file content](#result-markdown-file-content)

## Summary

...
```
