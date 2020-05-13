# automated cluster deployment

## what this is about
To pack all required tools for automated deployment of a cluster infrastructure, this repository was created. It includes the generation of ssh keys, a pxe server in dns-proxy mode and an ansible server combined in a container.
It can/will be decided later in the process whether the pxe-server and/or the ansible server will be migrated into the cluster.
The key generation is optional at this point, as by default a public ssh key is loaded from a provided url.

## how to use

### prerequisites
- install docker
- install go-task (https://taskfile.dev)
- run ```task init```
- run ```task run```
- ensure the 'to-be-pxe-booted' host's BIOS is set to boot from network first. Debian will add itself to bootmenu too during installation. Afterwards recheck the bootmenu and make sure network is still number one and debian listed second. Disable every other bootdevice. This may be done automatically in the future (#TODO).


- update readme with proper howto/troubleshooting chapters
  - make sure boot order has network before debian/disk, but have both enabled, so if network exits, the boot from disk runs
  - consolidate variables for whole project
  - consolidate variables/files for each machine

## notes
Most reoccuring tasks are specified in the ```./Taskfile.yml```.
If further information is needed, please look at ```./pxe/readme.md```.

## troubleshooting

### timings
To get a feeling for when a problem occured, here are some measured times for a bare-metal-deployment of a single machine:
- from pressing the boot button until accessible: ~ 7-15min
- poweroff: ~ 5s
- (normal) boot: ~ 40s
- network, then normal boot (pxe online): ~ 1:15
- network, then normal boot (pxe offline): ~ 1:40 (pxe timeout is a very long time... probably a minute)

### boot priority changes via cli
- only for next boot:
  efibootmgr --bootnext 00XX
  efibootmgr -n 00XX
- oneliner, to deactivate debian entry:
  sudo efibootmgr -A -b $(efibootmgr | grep "^Boot....\* debian$" | cut -c5-8)
  sudo efibootmgr --inactive -b $(efibootmgr | grep "^Boot....\* debian$" | cut -c5-8)
- oneliner, to activate debian entry:
  sudo efibootmgr -a -b $(efibootmgr | grep "^Boot....  debian$" | cut -c5-8)
  sudo efibootmgr --active -b $(efibootmgr | grep "^Boot....  debian$" | cut -c5-8)
- to list the current settings
  efibootmgr
  -> BootOrder: 0003,0001,0002
- to change settings (persistent)
  bootmgr -o 0002,0001,0003
- to change settings (single reboot)
  efibootmgr --bootnext 0002
