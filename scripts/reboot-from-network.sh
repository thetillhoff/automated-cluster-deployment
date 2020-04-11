#!/bin/bash

read -p "Please enter the host you want to reboot from network: " host
read -p "Please enter the user you want to be: [enforge] " user
user=${user:=enforge}  # if variable not set or null, set it to default

read -p "Are you SURE, you don't need to just run a reboot? [y|n] " -n 1 -r
echo # to move cursor to next line
if [[ $REPLY =~ ^[Yy]$ ]]
then
  ssh -t -i ../key -o "StrictHostKeyChecking=no" -o "UserKnownHostsFile=/dev/null" -o "LogLevel=error" "$user@$host" 'sudo efibootmgr -A -b $(efibootmgr | grep "^Boot....\* debian$" | cut -c5-8) && sudo systemctl reboot'
fi
echo # start new line
