# `toc`

## Table of contents

- [Summary](#summary)
- [Manual](#manual)
- [Usage](#usage)
- [Example](#example)
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
              e.g. toc --document=README.md --header="Table of contents" --limiter-right="##" int.
   ext      Generate table of content within a document and use provided paths as a list 
              e.g. toc --document=README.md --header="Table of contents" --limiter-right="##" ext --path=one/document.md --path=two/document.md.
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --document value       markdown file path to be updated
   --header value         header to use for document lookups and generation
   --limiter-left value   string to use as a lookup limiter (default: "##")
   --limiter-right value  string to use as a lookup limiter - empty will use end of file as a limit (default: "##")
   --help, -h             show help
```

## Usage

## Example

### Original Markdown file content

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

### Command

The following command will take all provided paths and generate a list of links within a document.

``` shell
notatio toc --document=README.md --header="Table of contents" --limiter-right=## ext --path=one/one.md --path=two/two.md
```

### Result Markdown file content

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
