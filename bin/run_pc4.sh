#!/bin/bash

export MUSIC_FOLDER="/mnt/nasmedia"
export KNOWN_ALBUMS_FILE="/mnt/dashboard/knownalbums.txt"
export KNOWN_ALBUMS_BACKUP="/mnt/dashboard/knownalbums_backup.txt"
export LOG_FILE="/home/ed/porthole.log"
export FETCH_INTERVAL=30000
export DASHBOARD_REFRESH_INTERVAL=5000
export SLEEP_AFTER=60000
export FOLDERS_TO_SCAN="Music/flac:3,Music/flac-cd:3,Music/flac-add:2,Music/flac-vorbis320:2,Music/mp3/main:2,Music_other/hd audio:3"
export LATEST_ADDITIONS_LIMIT=50
export GIT_COMMIT=`git log --pretty=format:'%h' -n 1`

echo "Mounting media folder from nas"
rm /home/ed/porthole.log
touch /home/ed/porthole.log
sudo mkdir -p /mnt/nasmedia
sudo mkdir -p /mnt/dashboard
sudo chmod 777 /mnt/nasmedia
sudo chmod 777 /mnt/dashboard
sudo sh -c 'mount -t cifs //192.168.1.102/media /mnt/nasmedia --verbose -o credentials=/home/ed/nascredentials >> /home/ed/porthole.log 2>&1'
sudo sh -c 'mount -t cifs //192.168.1.102/dashboard /mnt/dashboard --verbose -o credentials=/home/ed/nascredentials >> /home/ed/porthole.log 2>&1'

echo "Starting porthole in background..."
/home/ed/go/bin/porthole >> /home/ed/porthole.log 2>&1 &