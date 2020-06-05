#!/usr/bin/env php
<?php

include_once __DIR__ . DIRECTORY_SEPARATOR . '_utils.php';

getZDCmd();

if (!function_exists('simplexml_load_file')) {
    print("simplexml_load_file missing, pls use 'sudo apt-get install php7.0-simplexml' to install\n");
}