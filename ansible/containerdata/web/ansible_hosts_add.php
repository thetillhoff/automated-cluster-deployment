<?php
    header("Content-type: text/plain");


    ### php variables

    $clientip = $_SERVER['REMOTE_ADDR'];
    echo $clientip;

    
    ### variables and initial file creation

    $ansible_hosts_file = '/etc/ansible/hosts';

    $content = file_get_contents($ansible_hosts_file);
    $lines = explode("\n",$content);

    if(! in_array($clientip,$lines)) {
        $index = array_search("[nogroup]",$lines);
        $index = $index;
        
        $beginning = array_splice($lines,0,sizeof($lines)-$index);
        $newlines = array_merge($beginning, array($clientip));
        $lines = array_merge($newlines, $lines);
    }

    $content = implode("\n",$lines);

    $handle = fopen($ansible_hosts_file, 'w') or die('Cannot open file:  '.$ansible_hosts_file);
    fwrite($handle, $content);
    fclose($handle);
?>