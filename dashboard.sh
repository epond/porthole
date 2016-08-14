#!/bin/bash

export GOPATH="/home/pi/go"
export MUSIC_FOLDER="/mnt/nasmedia/Music"
export KNOWN_RELEASES_FILE="~/knownreleases"

echo "Waiting for network..."
/home/pi/go/src/github.com/epond/porthole/waitforip.sh

echo "Updating porthole..."
cd /home/pi/go/src/github.com/epond/porthole
rm /home/pi/porthole.log
git pull -r >> /home/pi/porthole.log 2>&1
bash -c ./build.sh >> /home/pi/porthole.log 2>&1

echo "Mounting media folder from nas"
sudo mkdir -p /mnt/nasmedia
sudo chmod 777 /mnt/nasmedia
sudo mount -o nolock 192.168.1.102:/volume1/media /mnt/nasmedia

echo "Starting porthole in background..."
/home/pi/go/bin/porthole >> /home/pi/porthole.log 2>&1 &

echo "Starting browser in kiosk mode..."
# https://github.com/elalemanyo/raspberry-pi-kiosk-screen
sudo -u pi /usr/bin/epiphany-browser -a -i --profile ~/.config http://localhost:9000 --display=:0 &
sleep 10s;
xte "key F11" -x:0
