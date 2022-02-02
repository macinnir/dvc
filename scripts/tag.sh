#!/bin/sh 

# Requires an initial annotated tag be set for the project 
# git tag -a "0.1.0" -m "0.1.0"
changes=`git rev-list $(git describe --abbrev=0)..HEAD --count`

if [ "$changes" = "0" ]; then 
    echo "No changes since last tag. Quitting..."
    exit 0 
fi 

branch=`git rev-parse --abbrev-ref HEAD`
versionType="${1}"

if [ "$branch" = "develop" ]; then 
    versionType="beta"
fi 

if [ "$branch" = "release" ]; then 
    versionType="rc"
fi 

newVersion=`dvc version $(git describe --tags) ${versionType}`

echo "Tagging with new ${versionType} version ${newVersion}"

# git tag -s # GPG signing 
git tag -a -m "New Tag: ${newVersion}" $newVersion
git push origin ${branch}
git push origin --tags
