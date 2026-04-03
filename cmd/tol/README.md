# `tol`

## Table of contents

- [Summary](#summary)
- [Manual](#manual)
- [Usage](#usage)

## Summary

A simple tool to generate a table of licenses from [`go-licenses`](https://github.com/google/go-licenses) output.

## Manual

``` bash
$ notatio tol --help
NAME:
   notatio tol

USAGE:
   notatio tol [command options]

DESCRIPTION:
   Table of licences generator

OPTIONS:
   --csv-path value       input csv file (go-licenses output)
   --document-path value  markdown document file path to be updated
   --header value         header to use for document lookups and generation
   --limiter-left value   string to use as a lookup limiter (default: "##")
   --limiter-right value  string to use as a lookup limiter - empty will use end of file as a limit (default: "##")
   --index value          index of a section to be used as a placeholder (useful if limiters refer to more than one section,
      0 = replace all) (default: 0)
   --skip value [ --skip value ]  packages to skip
   --help, -h                     show help
```

## Usage

You need the `go-licenses` to be installed first (please follow the instructions found on the
<https://github.com/google/go-licenses> page).

Run the following first to generate a **CSV** file.

``` bash
go-licenses report ./... > tmp/licenses.csv
```

Use the below **README.md** example.

``` markdown
# Title

## Licenses

| Package                                 | Licence                                                         | Type         |
|-----------------------------------------|-----------------------------------------------------------------|--------------|
| github.com/cpuguy83/go-md2man/v2/md2man | https://github.com/cpuguy83/go-md2man/blob/v2.0.7/LICENSE.md    | MIT          |
| github.com/russross/blackfriday/v2      | https://github.com/russross/blackfriday/blob/v2.1.0/LICENSE.txt | BSD-2-Clause |
| github.com/urfave/cli/v2                | https://github.com/urfave/cli/blob/v2.27.7/LICENSE              | MIT          |
| github.com/xrash/smetrics               | https://github.com/xrash/smetrics/blob/686a1a2994c1/LICENSE     | MIT          |

## Example
```

Follow up with the command below.

``` bash
notatio tol \
  --document-path=README.md \
  --csv-path=tmp/licenses.csv \
  --header="Licenses" \
  --limiter-left="##" \
  --limiter-right="## Example" \
  --index=1
```
