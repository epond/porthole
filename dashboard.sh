#!/bin/bash

echo "Waiting for network..."
/home/pi/go/src/github.com/epond/porthole/waitforip.sh

echo "Starting porthole in background..."
/home/pi/go/bin/porthole &

echo "Starting browser in kiosk mode..."
sudo -u pi /usr/bin/epiphany-browser -a -i --profile ~/.config http://localhost:9000 --display=:0 &
sleep 10s;
xte "key F11" -x:0
