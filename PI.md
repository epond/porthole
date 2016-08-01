# Enable SSH remote administration

    sudo raspi-config

# Configure the 7" touchscreen display

to rotate screen, edit `/boot/config.txt` and add the line: `lcd_rotate=2` to the top

reduce brightness (after enabling ssh to be safe):

    sudo su
    echo 80 > /sys/class/backlight/rpi_backlight/brightness
    exit

An introduction to using the display can be found here:
http://learn.pimoroni.com/tutorial/pi-lcd/getting-started-with-raspberry-pi-7-touchscreen-lcd

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
