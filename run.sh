#!/usr/bin/env bash
export MUSIC_FOLDER="/Users/edpond/Documents/porthole/dummymusic"
export KNOWN_RELEASES_FILE="/Users/edpond/Documents/porthole/knownreleases"
export GIT_COMMIT=`git log --pretty=format:'%h' -n 1`
export LOG_FILE="/Users/edpond/Documents/porthole/porthole.log"
export FETCH_INTERVAL=5
export DASHBOARD_REFRESH_INTERVAL=5
export FOLDERS_TO_SCAN=flac-vorbis320:2,mp3/main:2

./build.sh &&
porthole