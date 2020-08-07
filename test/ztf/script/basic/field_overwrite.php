#!/usr/bin/env php
<?php
/**
[case]
title=system、custom、default、config配置4层覆写
cid=1352
pid=7

[group]
 显示10行生成的数据 >>
 验证第3行数据    >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("default.yaml", "test.yaml", 10, "", array("fields"=>"field_nested_instant"));

$count = sprintf("%d", count($output));
print(">> $count\n");

print(">> $output[2]\n");