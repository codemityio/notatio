#!/bin/bash

# Note: the set [-/+] x is purely there to turn on and off outputting of the commands being executed.
if [ "${DEBUG}" = "true" ]; then
  set -x
fi

case "$1" in

"uml")
  echo "Coming soon..."
  ;;

"depgraph")
  echo "Coming soon..."
  ;;

"pkg")
  echo "Coming soon..."
  ;;

"render")
  echo "Coming soon..."
  ;;

"main")
  echo "Coming soon..."
  ;;

*)
  echo "error: incorrect '$1' command..."
  ;;

esac
