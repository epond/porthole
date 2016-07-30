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

# Install golang (currently 1.3.3)

    sudo apt-get install golang-go
    mkdir -p go/src/github.com/epond

add to `~/.bashrc`:

    export GOPATH="$HOME/go"
    export PATH="$PATH:$GOPATH/bin"
    
# Clone the porthole repo

    cd go/src/github.com/epond
    git clone https://github.com/epond/porthole.git

# Modify the PI's autostart script

replace the contents of `~/.config/lxsession/LXDE-pi/autostart` with:

    @xset s off
    @xset -dpms
    @xset s noblank
    @/home/pi/go/src/github.com/epond/porthole/dashboard.sh
