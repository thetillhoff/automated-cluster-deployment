<?php
    header("Content-type: text/plain");


    ### php variables

    $clientip = $_SERVER['REMOTE_ADDR'];
    echo $clientip."\r\n";

    
    ### variables and initial file creation

    $ansible_hosts_file = '/etc/ansible/hosts';

    $content = file_get_contents($ansible_hosts_file);
    $lines = explode("\n",$content);

    if(! in_array($clientip,$lines)) {
        $index_header = array_search("[nogroup]",$lines); # search for line [nogroup]
        $beginning = array_splice($lines,0,sizeof($lines)-$index_header); # detach everything from [nogroup] upwards (including [nogroup])

        $index_empty = array_search("",$lines); # search for first empty line after [nogroup]
        $second_beginning = array_splice($lines,0,sizeof($lines)-$index_empty-1); # detach everything from first empty line upwards (excluding empty line)
        $beginning = array_merge($beginning, $second_beginning); # merging everything before newlines into beginning

        $newlines = array_merge($beginning, array($clientip)); # appending the new line to beginning and put result in newlines
        $lines = array_merge($newlines, $lines); # merge newlines and lines

        $content = implode("\n",$lines); # merge array of lines with linebreaks

        $handle = fopen($ansible_hosts_file, 'w') or die('Cannot open file:  '.$ansible_hosts_file); # open file with 'w' permission
        fwrite($handle, $content); # write to file
        fclose($handle); # close filehandler
    }
?>