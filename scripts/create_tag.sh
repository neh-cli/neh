#!/usr/bin/env bash

# Usage: ./create_tag.sh <tag>
# Creates a new release and pushes it to the remote repository with the specified tag.
# Example: ./create_tag.sh v1.0.0

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <tag>"
  exit 1
fi

TAG=$1

./scripts/delete_tag.sh $TAG

# Create and push the new tag
git tag $TAG
git push origin $TAG

echo "Release $TAG has been created."
