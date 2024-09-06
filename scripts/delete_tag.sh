#!/usr/bin/env bash

# Usage: ./delete_tag.sh <tag>
# Deletes a specified Git tag and its corresponding GitHub release.
# Example: ./delete_tag.sh v1.0.0

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <tag>"
  exit 1
fi

TAG=$1

# Check if the tag exists
if git rev-parse "$TAG" >/dev/null 2>&1; then
  # Delete the tag
  git tag -d $TAG

  # Delete the remote tag
  if git push origin --delete $TAG 2>/dev/null; then
    echo "Remote tag $TAG has been deleted."
  else
    echo "Remote tag $TAG does not exist."
  fi

  # Delete the release
  if gh release view $TAG >/dev/null 2>&1; then
    gh release delete $TAG --yes
    echo "Tag and release $TAG have been deleted."
  else
    echo "No release corresponding to tag $TAG exists."
  fi

else
  echo "Tag $TAG does not exist. Skipping deletion."
fi
