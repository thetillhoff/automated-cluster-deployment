<?php
    header("Content-type: text/plain");

    ### php variables
    $serverip = $_SERVER['SERVER_ADDR'];
    $mac = str_replace(":", "-", htmlspecialchars($_GET["mac"]));
    $username = $_GET["username"];
    $keyboardlayout = "de";

?>#!/bin/bash

# enforce proper keyboard layout
echo \'# KEYBOARD CONFIGURATION files

# Consult the keyboard(5) manual page.

XKBMODEL="pc105"
XKBLAYOUT="<?php echo "$keyboardlayout"; ?>"
XKBVARIANT=""
XKBOPTIONS=""

BACKSPACE="guess"
\' > /etc/default/keyboard


# create firstboot.service file
echo '[Unit]
Description =   Run script at first startup after all services are loaded
Wants =         network-online.target
After =         network.target network-online.target
#                ^->former for compatibility, later actual target

[Service]
Type =          simple
ExecStartPre =  /bin/sleep 10
ExecStart =     /firstboot.sh

[Install]
WantedBy =      default.target
' > /etc/systemd/system/firstboot.service


echo '#!/bin/bash

# make sure permissions are correct on ~/.ssh folder
chown -R <?php echo "$username".":"."$username"; ?> /home/<?php echo "$username"; ?>/.ssh

# add self to ansible hosts list of pxe server
curl "http://<?php echo "$serverip"; ?>/ansible_hosts_add.php?mac=<?php echo "$mac"; ?>"

# leave ready-note
touch /home/<?php echo "$username"; ?>/ansible_ready

# remove firstboot files
systemctl disable firstboot
rm /etc/systemd/system/firstboot.service
rm /firstboot.sh
' > /firstboot.sh


chmod +x /firstboot.sh
systemctl enable firstboot
