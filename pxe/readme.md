# what this is about
This container will run a pxe proxy server, which works along with another primarey dhcp server and only adds up the pxe information to client which requested dhcp-information via broadcast. Boot data is passed via http too, so this container also packs an nginx on port 80.

## troubleshooting
If you have problems with the preseed-commands, note, that it's log are within the las ~30 lines of /var/log/installer/syslog after the installation is finished.

If you get the error "dnsmasq: failed to bind DHCP server socket: Address already in use" and if
  sudo netstat -tulpn | grep :53
returns systemd-resolve or dnsmasq and
  sudo netstat -tulpn | grep :67 
returns dnsmasq as occupying process, kill them via
  sudo kill <pid>
where <pid> is the number before the program name from the netstat command.
Restart them later (for internet access) via
  sudo service network-manager restart
