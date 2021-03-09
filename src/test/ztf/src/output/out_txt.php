#!/usr/bin/env php
<?php
/**
[case]
title=输出为Text格式
cid=1355
pid=7

[group]
 找到预期的文本1 >>
 找到预期的文本2 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 10, "output/default.txt", array("fields"=>"field_common"));

$lineArr = $zd->readOutput("output/default.txt");

$count = sprintf("%d", count($lineArr));
print(">> $count\n");

print(">> $lineArr[2]\n");