<?php

function getZDCmd()
{
    $os = strtolower(PHP_OS);
    $is64bit = is64bit();
    print("$os $is64bit \n");

    $ret = '';
    if ($is64bit && $os == 'darwin') {
        $ret = './zd-mac';
    } else if ($is64bit && $os == 'linux') {
        $ret = './zd-linux';
    } else if ($is64bit && strpos($os,"win") > -1) {
        $ret = 'zd-amd64.exe';
    } else if (!$is64bit && strpos($os,"win") > -1) {
        $ret = 'zd-x86.exe';
    }

    if ($ret == '') {
        die('Please test on 64/32 bits windows, or 64 bits linux, mac system.\n');
    } else {
        print("$ret \n");
    }

    return $ret;
}

function is64bit() {
    $int = "9223372036854775807";
    $int = intval($int);
    if ($int == 9223372036854775807) {
        /* 64bit */
        return true;
    } else if ($int == 2147483647) {
        /* 32bit */
        return false;
    } else {
        /* error */
        return false;
    }
}
