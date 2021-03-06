#!/usr/bin/env bash
export MUSIC_FOLDER="$GOPATH/src/github.com/epond/porthole"
export KNOWN_ALBUMS_FILE="$GOPATH/src/github.com/epond/porthole/knownalbums.txt"
export KNOWN_ALBUMS_BACKUP=$MUSIC_FOLDER"/knownalbums_backup.txt"
export GIT_COMMIT=`git log --pretty=format:'%h' -n 1`
export LOG_FILE="$GOPATH/src/github.com/epond/porthole/porthole.log"
export FETCH_INTERVAL=1000
export DASHBOARD_REFRESH_INTERVAL=1000
export SLEEP_AFTER=1000
export LATEST_ADDITIONS_LIMIT=10
export FOLDERS_TO_SCAN="testdata:2"

./build.sh &&
$GOPATH/bin/porthole
