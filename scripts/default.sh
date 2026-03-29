#!/bin/bash

if ! command -v auxilium &>/dev/null; then
  echo "Please install (https://github.com/${VENDOR}/auxilium) tools..."
else
  list=$(grep -E '^[a-zA-Z0-9_-]+:.*?## .*$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%s %s\n", $1, $2}')
  targets=$(auxilium select --select-name-label="Target" --select-value-label="Description" --list="$list" ${SIZE:+--size=${SIZE}})
  read -r -a parts <<<"$targets"
  if [[ ${parts[0]} == "" ]]; then exit 0; fi
  # shellcheck disable=SC2128
  $1 "${parts}"
fi
