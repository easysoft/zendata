#!/usr/bin/env php
<?php
/**
[case]
title=內置Instances一层引用
cid=1370
pid=7

[group]
 显示10行生成的数据 >>
 验证第3行数据     >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "test.yaml", 10, "", array("fields"=>"field_nested_instant"));

$count = sprintf("%d", count($output));
print(">> $count\n");

print(">> $output[2]\n");