#!/bin/bash

## create /container/tftp/boot.ipxe with current ip
# note that initrd is used as "ifexists"
# if multiple macs exists maybe ${net0/mac:hexhyp} should be used
echo -e "#!ipxe\n" > /container/tftp/boot.ipxe
echo -e "set boot-url http://$(cut -d' ' -f1 <<<$(hostname -I))/\n" >> /container/tftp/boot.ipxe
echo -e "initrd \${boot-url}mac/\${mac:hexhyp}.lastinstalled && chain \${boot-url}boot_harddrive.php?mac=\${mac:hexhyp} ||" >> /container/tftp/boot.ipxe # if already installed, boot from harddrive
echo -e "chain \${boot-url}mac/\${mac:hexhyp}.php?mac=\${mac:hexhyp} ||" >> /container/tftp/boot.ipxe # if mac is known and config file exists, load it
echo -e "chain \${boot-url}default.php?mac=\${mac:hexhyp} ||" >> /container/tftp/boot.ipxe # else load menu
echo -e "shell" # else drop to shell

## set ip in dnsmasq.conf
#sed -i "s/dhcp-boot=tag:ipxe,http:\/\/someip\/boot.ipxe/dhcp-boot=tag:ipxe,http:\/\/$(cut -d' ' -f1 <<<$(hostname -I))\/boot.ipxe/g" /etc/dnsmasq.conf

## allow nginx to create files in /container/web/mac folder
chmod -R o+w /container/web/mac
## allow nginx to read files from /key folder
chmod -R o+r /key
## allow nginx to write to ansible-hosts file
chmod o+w /etc/ansible/hosts
