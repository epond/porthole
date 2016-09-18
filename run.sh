#!/usr/bin/env bash
export MUSIC_FOLDER="/Users/edpond/Documents/porthole/dummymusic"
export KNOWN_RELEASES_FILE="/Users/edpond/Documents/porthole/knownreleases.txt"
export KNOWN_RELEASES_BACKUP=$MUSIC_FOLDER"/knownreleases_backup.txt"
export GIT_COMMIT=`git log --pretty=format:'%h' -n 1`
export LOG_FILE="/Users/edpond/Documents/porthole/porthole.log"
export FETCH_INTERVAL=1
export DASHBOARD_REFRESH_INTERVAL=1
export LATEST_ADDITIONS_LIMIT=10
export FOLDERS_TO_SCAN=flac-vorbis320:2,mp3/main:2

./build.sh &&
porthole