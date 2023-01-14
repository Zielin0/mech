#!/bin/sh

set -e

if [ "$USER" != "root" ]
then
  echo "ERROR: Permission denied"
  exit 1
else
  {
    echo "INFO: Building mech..."
    go build
    echo "INFO: Success"

    echo ""

    {
      echo "INFO: Installing mech..."
      mv ./mech /usr/local/bin
      echo "INFO: Success"
    } || {
      echo "ERROR: Failed to install mech"
      exit 1
    }
  } || {
    echo "ERROR: Go compiler needed"
    exit 1
  }
fi