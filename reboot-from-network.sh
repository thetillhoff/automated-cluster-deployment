#!/bin/bash

read -p "Please enter the host you want to reboot from network: " host
read -p "Please enter the user you want to be: [enforge] " user
user=${user:=enforge}  # if variable not set or null, set it to default

ssh -t -i ./key -o StrictHostKeyChecking=no "$user@$host" 'sudo efibootmgr -A -b $(efibootmgr | grep "^Boot....\* debian$" | cut -c5-8) && sudo systemctl reboot'
