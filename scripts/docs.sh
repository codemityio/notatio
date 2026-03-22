#!/bin/bash

# Note: the set [-/+] x is purely there to turn on and off outputting of the commands being executed.
if [ "${DEBUG}" = "true" ]; then
  set -x
fi

case "$1" in

"uml")
  if [ -z "${PACKAGES}" ]; then
    packages=$(find "pkg" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
    targets=$(auxilium select --label="Choose package to generate documents for" --select-name-label="Target package" --list="${packages}")
    if [[ ${targets} == "" ]]; then exit 0; fi
  else
    targets=${PACKAGES}
  fi
  for target in ${targets//,/ }; do
    echo "pkg/${target}/..."
    goforma code uml \
      --include-var \
      --include-const \
      --include-func \
      --include-not-exported \
      --workdir "${PWD}" \
      --json-output-path "pkg/${target}/graph.json" \
      --path "./pkg/${target}/..." >"pkg/${target}/graph.puml"
    docker run --rm \
      -w "${PWD}" \
      -v "${PWD}:${PWD}" \
      "${VENDOR}"/notatio:latest plantuml --input-path="pkg/${target}/graph.puml" --output-format=svg
  done
  ;;

"depgraph")
  if [ -z "${PACKAGES}" ]; then
    packages=$(find "pkg" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
    targets=$(auxilium select --label="Choose package to generate documents for" --select-name-label="Target package" --list="${packages}")
    if [[ ${targets} == "" ]]; then exit 0; fi
  else
    targets=${PACKAGES}
  fi
  for target in ${targets//,/ }; do
    echo "pkg/${target}/..."
    goforma code dep \
      --path "./pkg/${target}/..." \
      --workdir "${PWD}" \
      --exclude-standard \
      --exclude-vendor \
      --owned "${GOPRIVATE}" \
      >"pkg/${target}/depgraph.dot"
    docker run --rm \
      -v "${PWD}:${PWD}" \
      -w "${PWD}" \
      "${VENDOR}"/notatio:latest graphviz --input-path="pkg/${target}/depgraph.dot" --output-format=svg
  done
  ;;

"cmd")
  if [ -z "${PACKAGES}" ]; then
    packages=$(find "cmd" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
    targets=$(auxilium select --label="Choose package to generate documents for" --select-name-label="Target package" --list="${packages}")
    if [[ ${targets} == "" ]]; then exit 0; fi
  else
    targets=${PACKAGES}
  fi
  for target in ${targets//,/ }; do
    echo "cmd/${target}/..."
    notatio coi --command="${BASE_NAME} ${target} --help" --document="cmd/${target}/README.md" --header=Manual --limiter-left=## --limiter-right="## Usage"
    notatio toc --document="cmd/${target}/README.md" --header="Table of contents" --limiter-left="##" --limiter-right="## Summary" \
      int --start-from-level=1 --start-from-item=1
    docker run --rm \
      -v "${PWD}:${PWD}" \
      -w "${PWD}" \
      "${VENDOR}"/pandoc \
      --wrap=auto --columns=120 \
      --from=markdown-implicit_figures \
      --to=gfm --output="cmd/${target}/README.md" "cmd/${target}/README.md"
  done
  ;;

"pkg")
  if [ -z "${PACKAGES}" ]; then
    packages=$(find "pkg" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
    targets=$(auxilium select --label="Choose package to generate documents for" --select-name-label="Target package" --list="${packages}")
    if [[ ${targets} == "" ]]; then exit 0; fi
  else
    targets=${PACKAGES}
  fi
  for target in ${targets//,/ }; do
    echo "pkg/${target}/..."
    notatio toc --document="pkg/${target}/README.md" --header="Table of contents" --limiter-left="##" --limiter-right="## Summary" \
      int --start-from-level=1 --start-from-item=1
    docker run --rm \
      -v "${PWD}:${PWD}" \
      -w "${PWD}" \
      "${VENDOR}"/pandoc \
      --wrap=auto --columns=120 \
      --from=markdown-implicit_figures \
      --to=gfm --output="pkg/${target}/README.md" "pkg/${target}/README.md"
  done
  ;;

"render")
  if [ -z "${PACKAGES}" ]; then
    packages=$(find "pkg" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
    targets=$(auxilium select --label="Choose package to generate documents for" --select-name-label="Target package" --list="${packages}")
    if [[ ${targets} == "" ]]; then exit 0; fi
  else
    targets=${PACKAGES}
  fi
  for target in ${targets//,/ }; do
    echo "pkg/${target}/..."
    docker run --rm \
      -w "${PWD}" \
      -v "${PWD}:${PWD}" \
      "${VENDOR}"/notatio:latest mermaid --input-path="pkg/${target}" --output-format=svg --recursive
    docker run --rm \
      -w "${PWD}" \
      -v "${PWD}:${PWD}" \
      "${VENDOR}"/notatio:latest plantuml --input-path="pkg/${target}" --output-format=svg --recursive
    docker run --rm \
      -v "${PWD}:${PWD}" \
      -w "${PWD}" \
      "${VENDOR}"/notatio:latest graphviz --input-path="pkg/${target}" --output-format=svg --recursive
  done
  ;;

"main")
  # summaries
  packages=$(find "cmd" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
  paths=
  for target in ${packages//,/ }; do
    paths+=" --path=cmd/${target}/README.md"
  done
  notatio toc --document=README.md --header="Subcommands" --limiter-left="###" --limiter-right="###" \
    ext --summary-header="Summary" --summary-limiter-left="##" --summary-limiter-right="##" ${paths}
  packages=$(find "pkg" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
  paths=
  for target in ${packages//,/ }; do
    paths+=" --path=pkg/${target}/README.md"
  done
  notatio toc --document=README.md --header="Packages" --limiter-left="##" --limiter-right="##" \
    ext --summary-header="Summary" --summary-limiter-left="##" --summary-limiter-right="##" ${paths}
  # command
  notatio coi --command="${BASE_NAME} --help" --document=README.md --header=Manual --limiter-left=### --limiter-right=###
  # table of contents
  notatio toc --document=README.md --header="Table of contents" --limiter-right="## Summary" \
    int --start-from-level=1 --start-from-item=1
  docker run --rm \
    -v "${PWD}:${PWD}" \
    -w "${PWD}" \
    "${VENDOR}"/pandoc \
    --wrap=auto --columns=120 \
    --from=markdown-implicit_figures \
    --to=gfm --output=README.md README.md
  ;;

*)
  echo "error: incorrect '$1' command..."
  ;;

esac
