#!/bin/bash

export MUSIC_FOLDER="/mnt/nasmedia/Music"
export KNOWN_RELEASES_FILE="/mnt/dashboard/knownreleases.txt"
export KNOWN_RELEASES_BACKUP="/mnt/dashboard/knownreleases_backup.txt"
export LOG_FILE="/home/ed/porthole.log"
export FETCH_INTERVAL=30000
export DASHBOARD_REFRESH_INTERVAL=5000
export FOLDERS_TO_SCAN=flac:3,flac-cd:3,flac-add:2,flac-vorbis320:2,mp3/main:2
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