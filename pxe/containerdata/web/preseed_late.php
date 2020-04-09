<?php
    header("Content-type: text/plain");

    ### php variables
    $serverip = $_SERVER['SERVER_ADDR'];
    $mac = str_replace(":", "-", htmlspecialchars($_GET["mac"]));
    $username = $_GET["username"];

?>#!/bin/bash

echo '[Unit]
Description =   Run script at first startup after all services are loaded
After =         network.target network-online.target
#                ^->former for compatibility, later actual target
Wants =         network-online.target

[Service]
Type =          simple
ExecStart =     /firstboot.sh

[Install]
WantedBy =      default.target
' > /etc/systemd/system/firstboot.service

echo '#!/bin/bash

touch /home/enforge/test
chown -R <?php echo "$username".":"."$username"; ?> /home/<?php echo "$username"; ?>/.ssh
curl "http://<?php echo "$serverip"; ?>/ansible_hosts_add.php?mac=<?php echo "$mac"; ?>"

touch /home/enforge/test2

systemctl disable firstboot
rm /etc/systemd/system/firstboot.service
#rm /firstboot.sh
' > /firstboot.sh
chmod +x /firstboot.sh
systemctl enable firstboot
