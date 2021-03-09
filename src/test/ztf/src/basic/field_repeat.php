#!/usr/bin/env php
<?php
/**
[case]
title=重复字符{}
cid=1344
pid=7

[group]
 显示10行生成的数据 >>
 验证第4行数据    >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 10, "", array("fields"=>"field_repeat"));

$count = sprintf("%d", count($output));
print(">> $count\n");

print(">> $output[3]\n");