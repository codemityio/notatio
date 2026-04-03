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
