#!/usr/bin/env php
<?php
/**
[case]
title=输出为SQL格式
cid=1356
pid=7

[group]
 找到预期的SQL插入语句字段 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 10, "output/default.sql",
    array("fields"=>"field_common", "table"=>"tlb_table", "trim"=>"true"));

$lineArr = $zd->readOutput("output/field_common.sql");

print(">> $lineArr[2]\n");