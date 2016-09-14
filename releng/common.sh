#!/bin/bash

function confirm {
    read -r -p "${1} [y/N] " response
    case $response in
        [yY]) 
            true
            ;;
        *)
            false
            ;;
    esac
}

function panic_if_working_is_copy_dirty {
	if [ $(git status --porcelain | wc -l) -ne 0 ]; then
		echo "Error: Git working directory is dirty."
		exit 1
	fi
}

function execute {
	CMD=${1}
	echo -e "---> Executing: $CMD"
	eval $CMD
}

function error_exit
{
	echo "$1" 1>&2
	exit 1
}
