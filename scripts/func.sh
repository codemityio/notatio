#!/bin/bash

statan() {
  local path=$1
  local flag=$2
  export GOFLAGS="-buildvcs=false"
  eval "golangci-lint run -v ${flag} --timeout=2m ${path}"
  status=$?
  if [ ${status} -ne 0 ]; then
    echo "error: linter failed, please fix the errors..."
    exit ${status}
  fi
  eval "go vet ${path}"
  status=$?
  if [ ${status} -ne 0 ]; then
    echo "error: vetting failed, please fix the errors..."
    exit ${status}
  fi
  eval "govulncheck -show verbose ${path}"
  status=$?
  if [ ${status} -ne 0 ]; then
    echo "error: vulnerability check failed, please fix the errors..."
    exit ${status}
  fi
}

pushImage() {
  local image="$1"

  # strip tag
  local repo_with_namespace="${image%%:*}"

  # extract namespace and repo
  local namespace="${repo_with_namespace%%/*}"
  local repo="${repo_with_namespace#*/}"

  # check repository exists in Docker Hub
  if ! curl -fsSL "https://hub.docker.com/v2/repositories/${namespace}/${repo}/" >/dev/null 2>&1; then
    echo "error: Docker Hub repository does not exist: ${namespace}/${repo}"
    exit 1
  fi

  echo "pushing ${image}..."
  docker push "${image}"
}
