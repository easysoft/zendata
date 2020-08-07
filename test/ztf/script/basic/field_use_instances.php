#!/usr/bin/env php
<?php
/**
[case]
title=引用內置Instances
cid=1348
pid=7

[group]
 显示10行生成的数据 >>
 验证第3行数据    >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 10, "", array("fields"=>"field_use_instance"));

$count = sprintf("%d", count($output));
print(">> $count\n");

print(">> $output[2]\n");