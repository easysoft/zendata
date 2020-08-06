#!/usr/bin/env php
<?php
/**
[case]
title=
cid=0
pid=0

[group]
 >>
 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 10, "", array("fields"=>"field_file"));

$count = sprintf("%d", count($output));
print(">> $count\n");

$str = join($output, ",");
print("$str\n");

$array = array("aaron", "ben", "carl");
if (in_array($output[0], $array)) {
    print(">> $output[0] is in array\n");
}