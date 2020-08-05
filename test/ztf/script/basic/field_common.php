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

include_once __DIR__ . DIRECTORY_SEPARATOR . '../common/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 3, "", array("fields"=>"field_common"));

$count = sprintf("%d", count($output));
print(">> $count\n");
print(">> $output[0]\n");