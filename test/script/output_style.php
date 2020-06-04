#!/usr/bin/env php
<?php
/**
[case]
title=basic
cid=0
pid=0

[group]
  1. count of text >>
  2. value from text >>
  3. value from json >>
  4. value from xml >>
  5. first line of sql >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '_utils.php';
$cmd = getZDCmd();

$output = [];
exec("$cmd -y ../test/definition/basic.yaml -c 3 -field char,numb -o ../test/output/output.txt -f text", $output);
$content = file_get_contents('../test/output/output.txt');
$arr = explode("\n", $content);
$count = sprintf("%d", count($arr));
print(">> $count\n");
print(">> $arr[0]\n");

$output = [];
exec("$cmd -y ../test/definition/basic.yaml -c 3 -field char,numb -o ../test/output/output.json -f json", $output);
$content = file_get_contents('../test/output/output.json');
$json = json_decode($content);
$val = $json[0][0];
print(">> $val\n");

$output = [];
exec("$cmd -y ../test/definition/basic.yaml -c 3 -field char,numb -o ../test/output/output.sql -f sql", $output);
$content = file_get_contents('../test/output/output.sql');
$arr = explode("\n", $content);
print(">> $arr[0]\n");

if (function_exists('simplexml_load_file')) {
    print("lib simplexml_load_file missing, pls use 'sudo apt-get install php7.0-simplexml' to install\n");
} else {
    $output = [];
    exec("$cmd -y ../test/definition/basic.yaml -c 3 -field char,numb -o ../test/output/output.xml -f xml", $output);
    $xml = simplexml_load_file('../test/output/output.xml'); // sudo apt-get install php7.0-simplexml
    $val = $xml->table->row->col[0];
    print(">> $val\n");
}
