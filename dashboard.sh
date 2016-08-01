#!/bin/bash

echo "Waiting for network..."
/home/pi/go/src/github.com/epond/porthole/waitforip.sh

echo "Mounting media folder from nas"
sudo mkdir -p /mnt/nasmedia
sudo chmod 777 /mnt/nasmedia
sudo mount -o nolock 192.168.1.102:/volume1/media /mnt/nasmedia

echo "Updating porthole..."
cd /home/pi/go/src/github.com/epond/porthole
sudo -u pi git pull -r
export GOPATH="/home/pi/go"
bash -c ./build.sh

echo "Starting porthole in background..."
rm /home/pi/porthole.log
/home/pi/go/bin/porthole 2> /home/pi/porthole.log &

echo "Starting browser in kiosk mode..."
# https://github.com/elalemanyo/raspberry-pi-kiosk-screen
sudo -u pi /usr/bin/epiphany-browser -a -i --profile ~/.config http://localhost:9000 --display=:0 &
sleep 10s;
xte "key F11" -x:0
