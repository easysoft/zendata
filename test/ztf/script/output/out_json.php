#!/usr/bin/env php
<?php
/**
[case]
title=输出为JSON格式
cid=1357
pid=7

[group]
 找到预期的JSON字段 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 10, "output/default.json", array("fields"=>"field_common"));

$lineArr = $zd->readOutput("output/default.json");

print(">> $lineArr[2]\n");