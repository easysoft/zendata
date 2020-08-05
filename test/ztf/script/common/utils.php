<?php

include_once __DIR__ . DIRECTORY_SEPARATOR . 'config.php';

function getZdCmd()
{
    global $config;

    $os = strtolower(PHP_OS);
    $osBits = getOsBits();
    if (strpos($os, "win") == 0) {
        $os = $os . $osBits;
    } else if ($os == "darwin") {
        $os = "mac";
    }
    print("os = $os \n");

    $ret = $config['zd'][$os];

    if ($ret == '') {
        die('Please test on windows, linux or mac os.\n');
    } else {
        print("cmdPath is $ret \n");
    }

    return $ret;
}

function getDemoDir()
{
    $path = "../../demo";
     print("demo dir is $path\n");
    return $path;
}

function getOsBits() {
    $int = "9223372036854775807";
    $int = intval($int);
    if ($int == 9223372036854775807) {
        return "64";
    } else if ($int == 2147483647) {
        return "32";
    } else {
        return "";
    }
}
