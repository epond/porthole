# Getting started with the 7" touchscreen display

http://learn.pimoroni.com/tutorial/pi-lcd/getting-started-with-raspberry-pi-7-touchscreen-lcd

# Rotate screen

edit `/boot/config.txt` and add the line: `lcd_rotate=2` to the top

# Setup wifi

`sudo nano /etc/wpa_supplicant/wpa_supplicant.conf`

    network={
        ssid="The_ESSID_from_earlier"
        psk="Your_wifi_password"
    }

`sudo reboot`

# Install scala

download scala-2.11.8.tgz from http://www.scala-lang.org/download/

    sudo mkdir /usr/lib/scala
    sudo tar -xf scala-2.11.8.tgz -C /usr/lib/scala
    sudo ln -s /usr/lib/scala/scala-2.11.8/bin/scala /bin/scala
    sudo ln -s /usr/lib/scala/scala-2.11.8/bin/scalac /bin/scalac

# Install golang (currently 1.3.3)

    sudo apt-get install golang-go
    mkdir -p go/src/github.com/epond

add to ~/.bashrc:

    export GOPATH="$HOME/go"
    export PATH="$PATH:$GOPATH/bin"
