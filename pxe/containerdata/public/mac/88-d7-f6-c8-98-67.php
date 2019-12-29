#!ipxe

#hostname: BEAST
# The only variable thing in here is the actual filename and the hostname comment on line 3 - everything else is generic.

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


    ### variables

    $mac = str_replace(":", "-", htmlspecialchars($_GET["mac"]));


    ### last seen file creation
    
    $mac_file = '/container/public/mac/'.$mac.'_lasthandled';
    $handle = fopen($mac_file, 'w') or die('Cannot open file:  '.$mac_file);
    date_default_timezone_set('Europe/Berlin');
    $data = date('Y-m-d H:i:s', time());
    fwrite($handle, $data);
?>

### tell the user about the new mac address and sleep for X seconds
echo Existing MAC detected. File used from ./mac/<?php echo $mac; ?>.php
sleep 5

### ipxe variables

# base-url must end with a slash
set boot-url http://<?php echo $serverip; ?>/

# Figure out if client is 64-bit capable
cpuid --ext 29 && set arch x64 || set arch x86
cpuid --ext 29 && set archl amd64 || set archl i386


### operating systems

:debian-netinst
echo Booting Debian 10.0.0 netinstaller
kernel ${boot-url}netboot/debian-installer/amd64/linux initrd=initrd.gz
initrd ${boot-url}netboot/debian-installer/amd64/initrd.gz
#imgargs linux initrd=initrd.gz auto=true fb=false
# hostname and domain are set before preseed file is loaded, so they have to be set before
imgargs linux initrd=initrd.gz auto=true fb=false url=${boot-url}mac/<?php echo $mac; ?>_preseed.php hostname=BEAST domain=nodomain
boot || goto failed
goto start