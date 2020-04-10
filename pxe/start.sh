#!/bin/bash

## create /container/tftp/boot.ipxe with current ip
# note that initrd is used as "ifexists"
# if multiple macs exists maybe ${net0/mac:hexhyp} should be used
echo -e "#!ipxe\n" > /container/tftp/boot.ipxe
echo -e "set boot-url http://$(cut -d' ' -f1 <<<$(hostname -I))/\n" >> /container/tftp/boot.ipxe
#echo -e "initrd \${boot-url}mac/\${net0/mac:hexhyp}.lasthandled && chain \${boot-url}boot_harddrive.php?mac=\${net0/mac:hexhyp} ||" >> /container/tftp/boot.ipxe # if already installed, boot from harddrive
echo -e "initrd \${boot-url}mac/\${net0/mac:hexhyp}.lasthandled && goto bootnext ||" >> /container/tftp/boot.ipxe # if already installed, boot from harddrive
echo -e "chain \${boot-url}mac/\${net0/mac:hexhyp}.php?mac=\${net0/mac:hexhyp} ||" >> /container/tftp/boot.ipxe # if mac is known and config file exists, load it
echo -e "chain \${boot-url}default.php?mac=\${net0/mac:hexhyp} ||" >> /container/tftp/boot.ipxe # else load default menu
echo -e "shell" >> /container/tftp/boot.ipxe # else drop to shell
echo -e ":bootnext" >> /container/tftp/boot.ipxe # label for boot with next device (local disk)
echo -e "echo Booting from local harddrive in 1 seconds." >> /container/tftp/boot.ipxe # display message
echo -e "sleep 1" >> /container/tftp/boot.ipxe # give the user the opportunity to read the message
echo -e "exit" >> /container/tftp/boot.ipxe # exit and boot from next device (local disk)
## set ip in dnsmasq.conf
#sed -i "s/dhcp-boot=tag:ipxe,http:\/\/someip\/boot.ipxe/dhcp-boot=tag:ipxe,http:\/\/$(cut -d' ' -f1 <<<$(hostname -I))\/boot.ipxe/g" /etc/dnsmasq.conf

## allow nginx to create files in /container/web/mac folder
chmod -R o+w /container/web/mac
## allow nginx to read files from /key folder
chmod -R o+r /key
## allow nginx to write to ansible-hosts file
chmod o+w /etc/ansible/hosts
