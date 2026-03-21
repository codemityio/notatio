# `coi`

## Table of contents

- [Summary](#summary)
- [Manual](#manual)
- [Usage](#usage)

## Summary

A simple tool to generate document sections with provided command output.

## Manual

``` bash
$ notatio coi --help
NAME:
   notatio coi

USAGE:
   notatio coi [command options]

DESCRIPTION:
   Command output injector.

OPTIONS:
   --document value       markdown file path to be updated
   --header value         header to use for document lookups and generation
   --limiter-left value   string to use as a lookup limiter (default: "##")
   --limiter-right value  string to use as a lookup limiter - empty will use end of file as a limit (default: "##")
   --shell-name value     shell name to use in the output (default: "bash")
   --shell-prompt value   shell prompt prefix to use in the output (default: "$")
   --command value        command to execute (command execution is skipped if --output is also provided)
   --output value         output to inject
   --help, -h             show help
```

## Usage

Coming soon…
