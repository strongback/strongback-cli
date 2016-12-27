#!/usr/bin/env bash

ROOT_DIR=$(cd $(dirname $(dirname $0)) && pwd)

set -ex

component=$1

old_version=$(cat VERSION)
major=$(echo $old_version | cut -d'.' -f 1)
minor=$(echo $old_version | cut -d'.' -f 2)
patch=$(echo $old_version | cut -d'.' -f 3)

case "$component" in
  major )
    major=$(expr $major + 1)
    minor=0
    patch=0
    ;;
  minor )
    minor=$(expr $minor + 1)
    patch=0
    ;;
  patch )
    patch=$(expr $patch + 1)
    ;;
  * )
    echo "Error - argument must be 'major', 'minor' or 'patch'"
    echo "Usage: bump-version [major | minor | patch]"
    exit 1
    ;;
esac

version=$major.$minor.$patch

echo "Updating VERSION file to $version"
echo $version > VERSION

echo "Committing change"
git reset .
git add VERSION

git commit -m "Bump version to $version"

echo "Creating v$version tag"
git tag v$version

$ROOT_DIR/bin/generate-changelog $old_version $version
$ROOT_DIR/bin/commit-version-bump

echo -e "All Done! Review commit log and push to GitHub if all looks good."