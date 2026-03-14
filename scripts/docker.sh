#!/bin/bash

source scripts/func.sh

# Note: the set [-/+] x is purely there to turn on and off outputting of the commands being executed.
if [ "${DEBUG}" = "true" ]; then
  set -x
fi

case "$1" in

"tag")
  docker tag "${IMAGE_NAME}":latest "${IMAGE_NAME}:$(scripts/tools.sh version)"
  ;;

"push")
  pushImage "${IMAGE_NAME}:latest"
  pushImage "${IMAGE_NAME}:$(scripts/tools.sh version)"
  ;;

*)
  echo "error: incorrect '$1' command..."
  ;;

esac
