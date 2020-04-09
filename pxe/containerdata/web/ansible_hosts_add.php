<?php
    header("Content-type: text/plain");


    ### php variables

    $clientip = $_SERVER['REMOTE_ADDR'];
    echo $clientip;
    $mac = str_replace(":", "-", htmlspecialchars($_GET["mac"]));

    
    ### variables and initial file creation

    $ansible_hosts_file = '/etc/ansible/hosts';

    $content = file_get_contents($ansible_hosts_file); # get file contents
    $content = str_replace("\r\n","\n",$content); # ensure linebreaks are in the right format
    $lines = explode("\n",$content); # split on linebreaks into array

    if(! in_array($clientip,$lines)) { # if clientip not already contained
        $index_header = array_search("[nogroup]",$lines); # search for line [nogroup]
        $beginning = array_splice($lines,0,(sizeof($lines)-$index_header)-1); # detach everything from [nogroup] upwards (including [nogroup])

        $index_empty = array_search("",$lines); # search for first empty line after [nogroup]
        $second_beginning = array_splice($lines,0,(sizeof($lines)-$index_empty)-1); # detach everything from first empty line upwards (excluding empty line)
        $beginning = array_merge($beginning, $second_beginning); # merging everything before newlines into beginning

        $newlines = array_merge($beginning, array($clientip)); # appending the new line to beginning and put result in newlines
        $lines = array_merge($newlines, $lines); # merge newlines and lines

        $content = implode("\n",$lines); # merge array of lines with linebreaks

        $handle = fopen($ansible_hosts_file, 'w') or die('Cannot open file:  '.$ansible_hosts_file); # open file with 'w' permission
        fwrite($handle, $content); # write to file
        fclose($handle); # close filehandler

        echo " added.\n";
    } else {
        echo " is already listed.\n";
    }


    ### last seen file creation
    
    $mac_file = '/container/web/mac/'.$mac.'.lasthandled';
    $handle = fopen($mac_file, 'w') or die('Cannot open file:  '.$mac_file);
    date_default_timezone_set('Europe/Berlin');
    $data = date('Y-m-d H:i:s', time());
    fwrite($handle, $data);
?>