# open todos for this project

- update readme with proper troubleshooting chapter
  - make sure boot order has network before debian/disk, but have both enabled, so if network exits, the boot from disk runs
  - consolidate variables for whole project
  - consolidate variables/files for each machine
- edit efibootmgr settings automatically to reboot via network
  - only for next boot:
    efibootmgr --bootnext 00XX
    efibootmgr -n 00XX
  - oneliner, to deactivate debian entry:
    sudo efibootmgr -A -b $(efibootmgr | grep "^Boot....\* debian$" | cut -c5-8)
    sudo efibootmgr --inactive -b $(efibootmgr | grep "^Boot....\* debian$" | cut -c5-8)
  - onelines, to activate debian entry:
    sudo efibootmgr -a -b $(efibootmgr | grep "^Boot....  debian$" | cut -c5-8)
    sudo efibootmgr --active -b $(efibootmgr | grep "^Boot....  debian$" | cut -c5-8)
  - to list the current settings
    efibootmgr
    -> BootOrder: 0003,0001,0002
  - to change settings (persistent)
    bootmgr -o 0002,0001,0003
  - to change settings (single reboot)
    efibootmgr --bootnext 0002
- add other OS capabilites, like plain ubuntu and windows (10 & Server2019) (best would be iso-support)
- add mgmt GUI (e.g. via browser)
  - add reboot capabilities to that GUI
  - add different states of machines
    - offline
    - ipxe
    - installing
    - rebooting
    - online
  - add wake-on-lan capabilities to that GUI
