#!/bin/bash

# - See https://developer.github.com/guides/getting-started/ for
#   GitHub API authentication hints.
#
# - See https://developer.github.com/v3/repos/releases/ for 
#   GitHub API doc.

set -e
#set -x

function usage {
    echo "Usage: $0 [VERSION] [FILE] [GITHUB USERNAME]"
    echo "Example: $0 0.2.1 cntinsight-v0.2.2-x64-libmusl.tar.gz zaunerc"
    exit 1
}

if [ "$#" -ne 3 ]; then
	usage
fi

VERSION="$1"
FILE=$2
USER=$3

RELEASE_ID=$(curl -s https://api.github.com/repos/zaunerc/cntrinfod/releases/tags/v$VERSION | jq -r '.id')

curl -i -u $USER \
     -H "Content-Type: application/gzip" \
     --data-binary $FILE \
     "https://uploads.github.com/repos/zaunerc/cntrinfod/releases/$RELEASE_ID/assets?name=$FILE"

