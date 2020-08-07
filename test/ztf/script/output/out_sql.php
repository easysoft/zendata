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
$output = $zd->create("", "test2.yaml", 10, "output/test2.sql",
    array("fields"=>"test0", "table"=>"tlb_table", "trim"=>"true"));

$lineArr = $zd->readOutput("output/test2.sql");

print(">> $lineArr[2]\n");