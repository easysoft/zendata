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
$output = $zd->create("", "default.yaml", 10, "output/default.txt", array("fields"=>"field_common"));

$lineArr = $zd->readOutput("output/default.txt");

$count = sprintf("%d", count($lineArr));
print(">> $count\n");

print(">> $lineArr[2]\n");