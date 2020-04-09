<?php
    header("Content-type: text/plain");

    ### php variables
    $serverip = $_SERVER['SERVER_ADDR'];
    $mac = str_replace(":", "-", htmlspecialchars($_GET["mac"]));
    $username = $_GET["username"];

?>#!/bin/sh

echo '[Unit]
Description=Run script at first startup after all services are loaded
After=default.target

[Service]
Type=simple
#RemainAfterExit=no means there is no need to stop it at a later point
RemainAfterExit=no
ExecStart=/firstboot.sh

[Install]
WantedBy=default.target
' > /etc/systemd/system/firstboot.service

echo '#!/bin/bash

touch /home/enforge/test
chown -R <?php echo "$username".":"."$username"; ?> /home/<?php echo "$username"; ?>/.ssh
curl http://<?php echo "$serverip"; ?>/ansible_hosts_add.php?mac=<?php echo "$mac"; ?>

touch /home/enforge/test2

systemctl disable firstboot
rm /etc/systemd/system/firstboot.service
#rm /firstboot.sh
' > /firstboot.sh
chmod +x /firstboot.sh
systemctl enable firstboot
