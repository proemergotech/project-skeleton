#!/bin/sh
# Check for CHANGELOG updates and fails when not found 
if [ -f "CHANGELOG.md" ]
then
    if [ "$(git diff --name-only origin/master | grep -c "CHANGELOG.md")" -ne 1 ]
    then
        echo "Changelog should be updated!"
        exit 1
    fi
fi