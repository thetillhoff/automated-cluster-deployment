#!ipxe

<?php
    header("Content-type: text/plain");


    ### php variables

    $serverip = $_SERVER['SERVER_ADDR'];


    ### comments

    # mac:hexhyp results f.e. in 52-54-00-12-34-56
    #ipxe> echo ${net0/mac:hexhyp}
    #52-54-00-12-34-56
    #ipxe> echo ${net0/mac:hexraw}
    #525400123456

    # best source so far: https://gist.github.com/robinsmidsrod/2234639


    ### variables and initial file creation

    $mac = str_replace(":", "-", htmlspecialchars($_GET["mac"]));
    $mac_file = '/container/web/mac/'.$mac;
    $handle = fopen($mac_file, 'w') or die('Cannot open file:  '.$mac_file);
    $data = htmlspecialchars($_GET["mac"]);
    fwrite($handle, $data);
?>

### tell the user about the new mac address and sleep for X milliseconds
echo New MAC detected. File created at ./mac/<?php echo $mac; ?>
sleep 5000

### ipxe variables

# base-url must end with a slash
set boot-url http://<?php echo $serverip; ?>/

# Figure out if client is 64-bit capable
cpuid --ext 29 && set arch x64 || set arch x86
cpuid --ext 29 && set archl amd64 || set archl i386

# Set DNS server
set dns 1.1.1.1

# menu timeout in milliseconds
set menu-timeout 30000


### main menu

:start
set menu-default exit
menu iPXE boot menu for ${hostname} on ${ip}
item --gap --             ------------------------- Operating systems ------------------------------
item --key d debian-netinst    Debian 10.0.0 netinst
item --key w win10             Windows 10
item --gap --             ------------------------- Advanced options -------------------------------
item --key c config       Configure settings
item shell                Drop to iPXE shell
item reboot               Reboot computer
item
item --key x exit         Exit iPXE and continue BIOS boot
choose --timeout ${menu-timeout} --default ${menu-default} selected || goto cancel
goto ${selected}

:cancel
echo You cancelled the menu, dropping you to a shell

:shell
echo Type 'exit' to get the back to the menu
shell
set menu-timeout 0
set submenu-timeout 0
goto start

:failed
echo Booting failed, dropping to shell
goto shell

:reboot
reboot

:exit
exit
#sanboot --no-describe --drive 0x80

:config
config
goto start

:back
set submenu-timeout 0
clear submenu-default
goto start


### operating systems

:debian-netinst
echo Booting Debian 10.0.0 netinstaller
kernel ${boot-url}netboot/debian-installer/amd64/linux initrd=initrd.gz
initrd ${boot-url}netboot/debian-installer/amd64/initrd.gz
imgargs linux initrd=initrd.gz auto=true fb=false
boot || goto failed
goto start

:win10
echo Booting Windows 10 installer
# don't forget to set mime types for all of the below files to the webserver (all 'text/plain')
# . for wimboot and bcd
# .sdi for boot.sdi
# .wim for boot.wim
# or
# .iso for Windows10_en.iso
initrd http://192.168.1.166/Windows10_en.iso
chain ${boot-url}memdisk iso raw || goto failed
## bios mode
#kernel ${boot-url}wimboot
#kernel ${boot-url}wimboot index=2   # select bootimage
# file injection into X:\Windows\System32\
#initrd winpeshl.ini     winpeshl.ini
#initrd startup.bat      startup.bat
#initrd ${boot-url}win10/boot/bcd                    BCD
#initrd ${boot-url}win10/boot/boot.sdi               boot.sdi
#initrd ${boot-url}win10/sources/boot.wim            boot.wim
#boot || goto failed
goto start