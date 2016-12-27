#!/usr/bin/env bash

echo "Adding CHANGELOG"
git add CHANGELOG.md
echo "Ammending commit with CHANGELOG update"
git ci --amend --no-edit

echo "Retagging"
git tag -d v$(cat VERSION)
git tag v$(cat VERSION)