#!/bin/bash

# Used in https://github.com/zaunerc/go-release-scripts/blob/master/createTarGzPackage.sh.

function assembleArtifact {
	mkdir -p "$ARCHIVE_ROOT/usr/local/bin"
	mkdir -p "$ARCHIVE_ROOT/usr/local/etc/cntrinfod"
	mkdir -p "$ARCHIVE_ROOT/usr/local/share/cntrinfod"

	cp ../cntrinfod "$ARCHIVE_ROOT/usr/local/bin"
	cp -r ../static_data/* "$ARCHIVE_ROOT/usr/local/share/cntrinfod"
}

