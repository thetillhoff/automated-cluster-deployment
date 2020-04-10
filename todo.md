# open todos for this project

- edit efibootmgr settings automatically to reboot via network
  - efibootmgr | grep "^Boot00..\* debian$" returns the debian line
  - oneliner, to deactivate debian entry:
    sudo efibootmgr -A -b $(efibootmgr | grep "^Boot....\* debian$" | cut -c5-8)
  - onelines, to activate debian entry:
    sudo efibootmgr -a -b $(efibootmgr | grep "^Boot....  debian$" | cut -c5-8)
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
