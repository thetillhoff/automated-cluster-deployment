<?php
    header("Content-type: text/plain");

    ### php variables
    $serverip = $_SERVER['SERVER_ADDR'];
    $mac = str_replace(":", "-", htmlspecialchars($_GET["mac"]));
    $username = $_GET["username"];
    $keyboardlayout = "de";

?>#!/bin/bash

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

# out of unknown reasons, sometimes the network interface is 'down' after installation
#   - a reboot fixes this, but how to reboot if you don't have access to it?
# -> correction:
#for netint in /sys/class/net/!(lo); do ip link set $(basename $netint) up; done
#dhclient -r
#dhclient
#sleep 1

# add self to ansible hosts list of pxe server and leave debug message on success and error
curl -s "http://<?php echo "$serverip"; ?>/ansible_hosts_add.php?mac=<?php echo "$mac"; ?>" > /home/<?php echo "$username"; ?>/ansible_ready 2> /home/<?php echo "$username"; ?>/error
# delete empty files afterwards
find /home/<?php echo "$username"; ?>/ -type f -empty -delete

# remove firstboot files
systemctl disable firstboot
rm /etc/systemd/system/firstboot.service
rm /firstboot.sh
' > /firstboot.sh


chmod +x /firstboot.sh
systemctl enable firstboot
