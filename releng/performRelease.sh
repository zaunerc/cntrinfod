#!/bin/bash

set -e
#set -x

function usage {
    echo "Usage: $0 [VERSION] [NEXT VERSION]"
    echo "Example: $0 0.2.0 0.2.1"
    exit 1
}

if [ "$#" -ne 2 ]; then
	usage
fi

VERSION="$1"
NEXT_VERSION="$2-SNAPSHOT"

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd "$SCRIPT_DIR/.."
echo "Changed current working dir to $(pwd)"

source $SCRIPT_DIR/common.sh

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "master" ]; then
	error_exit "Releases are only suppored on master branch. Aborting!"
fi

panic_if_working_is_copy_dirty

# Update version info
execute 'sed --in-place "s/app.Version = \".*\"/app.Version = \"$VERSION\"/g" cntrinfod.go'

cat << EOF

************************************************************
* VERIFY GIT DIFF
*

$(git diff --no-color --unified=0 | cat)

*
************************************************************

************************************************************
* VERIFY GIT USERNAME AND EMAIL
*

Git user name: $(git config --get user.name)
Git email: $(git config --get user.email)

*
************************************************************

The following steps will be executed if you proceed:

1. master> Commit the version number change. New version is $VERSION.
2. master> Create an annotated tag.
3. Checkout branch dev
4. dev> Merge master into dev branch.
5. dev> Update version number to $NEXT_VERSION.
6. dev> Commit the version number change.

EOF

confirm "Proceed?"

# master branch

execute 'git commit -a -m "Release version $VERSION."'
execute 'git tag --annotate "v$VERSION" --message "Release version $VERSION."'

# dev branch

execute 'git checkout dev'
execute 'git merge master'
# Update version info
execute 'sed --in-place "s/app.Version = \".*\"/app.Version = \"$NEXT_VERSION\"/g" cntrinfod.go'
execute 'git commit -a -m "Starting development iteration on $NEXT_VERSION."'

cat << EOF

YOU CAN NOW PUSH THE CHANGES TO ORIGIN. EXECUTE:
\$ git push --follow-tags origin master dev

EOF

