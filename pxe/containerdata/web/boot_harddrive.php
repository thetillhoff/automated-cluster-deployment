#!ipxe

<?php
    header("Content-type: text/plain");
?>

### tell user that boot via harddrive is done
echo Booting from local harddrive in 1 seconds.
sleep 1


### reboot from local harddrive
#exit
# alternatively use
sanboot --no-describe --drive 0x80