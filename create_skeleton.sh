#!/usr/bin/env bash
set -e

if [[ "$@" < 1 ]]; then
    echo -e "\nUsage:\n"
    echo -e "  ./create_skeleton <target_folder_path>\n"
    exit 1
fi

# Dliver Project Skeleton folder name values
skeleton_underscore=dliver_project_skeleton
skeleton_dash=dliver-project-skeleton

# Target project folder
target_dir=$( echo $1 | sed 's:/*$::')  # Remove trailing slashes

project_name=${target_dir##*/}  # Extract project name from project dir
project_name_underscore=${project_name//[^a-zA-Z0-9]/_}
project_name_dash=${project_name//[^a-zA-Z0-9]/-}

# Exclude
#  .git               - version control related files / git history of project skeleton
#  .idea              - IntelliJ Idea/Goland related project information
#  vendor             - all 3rd party dependencies
#  create_skeleton.sh - this script file
#  *_test.go          - all test files
rsync -av --exclude=.git --exclude=.idea --exclude=vendor --exclude=create_skeleton.sh --exclude=*_test.go . ${target_dir}

# Replace skeleton project text and folder names to the new one
grep -rl ${skeleton_underscore} ${target_dir}/ | xargs sed -i "s@$skeleton_underscore@$project_name_underscore@g"
grep -rl ${skeleton_dash} ${target_dir}/ | xargs sed -i "s@$skeleton_dash@$project_name_dash@g"
