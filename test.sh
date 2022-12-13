#!/bin/sh
debug="$1"

export XBSAPI_API_URL="http://127.0.0.1:8000/api/v1"

failfast() {
  if [ "$1" -ne "0" ]
  then
    printf "   FAILED: %s\n" "$2"
    exit "$1"
  else
    printf "   SUCCESS\n"
    if [ "$debug" = "true" ]
    then
      printf "   DEBUG: %s\n" "$2"
    fi
  fi
}

printf "No tests yet.\n"

exit 0
