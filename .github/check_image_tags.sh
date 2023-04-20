#!/bin/bash

## Determine if image tags already exist
## Do not allow overwrites
function tag_exists() {
    IMAGE_TAG_BASE=${1#*/}
    curl -s -f -lSL https://hub.docker.com/v2/repositories/${IMAGE_TAG_BASE}/tags/$2 &> /dev/null
}

if tag_exists $1 $2; then
    echo tag $2 already exists on $1, exiting.
    exit 1
fi
