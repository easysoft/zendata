#!/usr/bin/env php
<?php
/**
[case]
title=输出为XML格式
cid=1358
pid=7

[group]
 找到预期的row元素 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 10, "output/default.xml", array("fields"=>"field_common"));

$lineArr = $zd->readOutput("output/default.xml");

print(">> $lineArr[4]\n");