#!/bin/bash

set -e
#set -x

function usage {
    echo "Usage: $0 [VERSION] [ARCH] [STDLIB]"
    echo "Example: $0 0.2.0 x64 libmusl"
    exit 1
}


if [ "$#" -ne 3 ]; then
	usage
fi

VERSION="$1"
ARCH="$2"
STDLIB="$3"

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
source $SCRIPT_DIR/common.sh

cd "$SCRIPT_DIR/../"
echo "Changed current working dir to $(pwd)"

panic_if_working_is_copy_dirty

execute 'git checkout v$VERSION'
execute 'go test'
execute 'go build'

cd "$SCRIPT_DIR"
echo "Changed current working dir to $(pwd)"

ARCHIVE_ROOT="tmp"

if [ -d "$ARCHIVE_ROOT" ]; then
  rm -rf "$ARCHIVE_ROOT"
fi

mkdir "$ARCHIVE_ROOT"

mkdir -p "$ARCHIVE_ROOT/usr/local/bin"
mkdir -p "$ARCHIVE_ROOT/usr/local/etc/cntrinfod"
mkdir -p "$ARCHIVE_ROOT/usr/local/share/cntrinfod"

cp ../cntrinfod "$ARCHIVE_ROOT/usr/local/bin"
cp -r ../static_data/* "$ARCHIVE_ROOT/usr/local/share/cntrinfod"

ARCHIVE="cntinsight-v$VERSION-$ARCH-$STDLIB.tar.gz"
tar -czf "$ARCHIVE" -C "$ARCHIVE_ROOT" usr

cat << EOF
************************************************************
* $ARCHIVE successfully created:
*

$(tar -tvzf "$ARCHIVE")

*
************************************************************
EOF
