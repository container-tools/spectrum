#!/bin/bash

set -e

if [ "$#" -ne 1 ]; then
  echo "usage: $0 version"
  exit 1
fi

version=$1

if [[ $version != v* ]]; then
  echo "version must start with \"v\""
  exit 1
fi

set +e
git diff-index --quiet HEAD --
if [ "$?" -ne 0 ]; then
  echo "there are uncommitted changes: please commit all changes before releasing"
  exit 1
fi
set -e

echo "Releasing $version..."

git tag -a $1 -m "Release $version"
git push upstream $version
