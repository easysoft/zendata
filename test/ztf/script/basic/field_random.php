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
$output = $zd->create("", "default.yaml", 10, "", array("fields"=>"field_random"));

$count = sprintf("%d", count($output));
print(">> $count\n");

$str = join($output, ",");
print("$str\n");

if ($output[0] != 1 || $output[1] != 2 || $output[2] != 3) {
    print(">> not 1,2,3\n");
}