#!/bin/bash

## configure dnsmasq for local ip
sed -i "s%dhcp-boot=tag:ipxe,http://xxx.xxx.xxx.xxx/default%dhcp-boot=tag:ipxe,http://$(cut -d' ' -f1 <<<$(hostname -I))/default?mac=\${net0/mac:hexhyp}%" /etc/dnsmasq.conf

## allow nginx to read files from /key folder
chmod -R o+r /key
## allow nginx to write to ansible-hosts file
chmod o+w /etc/ansible/hosts
