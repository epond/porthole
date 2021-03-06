#!/bin/bash

export GOPATH="/home/pi/go"
export MUSIC_FOLDER="/mnt/nasmedia"
export KNOWN_ALBUMS_FILE="/mnt/dashboard/knownalbums.txt"
export KNOWN_ALBUMS_BACKUP="/mnt/dashboard/knownalbums_backup.txt"
export LOG_FILE="/home/pi/porthole.log"
export FETCH_INTERVAL=1000
export DASHBOARD_REFRESH_INTERVAL=5000
export SLEEP_AFTER=1000
export LATEST_ADDITIONS_LIMIT=200
export FOLDERS_TO_SCAN="Music/flac:3,Music/flac-cd:3,Music/flac-add:2,Music/flac-vorbis320:2,Music/mp3/main:2,Music_other/hd audio:3"

echo "Waiting for network..."
/home/pi/go/src/github.com/epond/porthole/bin/waitforip.sh

cd /home/pi/go/src/github.com/epond/porthole
rm /home/pi/porthole.log
export GIT_COMMIT=`git log --pretty=format:'%h' -n 1`
bash -c ./build.sh >> /home/pi/porthole.log 2>&1

echo "Mounting media folder from nas"
sudo mkdir -p /mnt/nasmedia
sudo mkdir -p /mnt/dashboard
sudo chmod 777 /mnt/nasmedia
sudo chmod 777 /mnt/dashboard
sudo sh -c 'mount -t cifs //192.168.0.102/media /mnt/nasmedia --verbose -o vers=2.0,credentials=/home/pi/nascredentials >> /home/pi/porthole.log 2>&1'
sudo sh -c 'mount -t cifs //192.168.0.102/dashboard /mnt/dashboard --verbose -o vers=2.0,credentials=/home/pi/nascredentials >> /home/pi/porthole.log 2>&1'

echo "Starting porthole in background..."
/home/pi/go/bin/porthole >> /home/pi/porthole.log 2>&1 &
