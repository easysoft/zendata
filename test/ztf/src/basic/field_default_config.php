#!/usr/bin/env php
<?php
/**
[case]
title=default和config搭配使用
cid=1351
pid=7

[group]
 显示10行生成的数据 >>
 验证第3行数据     >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("default.yaml", "test.yaml", 10, "");

$count = sprintf("%d", count($output));
print(">> $count\n");

print(">> $output[2]\n");