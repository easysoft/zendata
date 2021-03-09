#!/usr/bin/env php
<?php
/**
[case]
title=Loop使用区间
cid=1374
pid=7

[group]
 显示10行生成的数据    >>
 验证第3行数据重复了4次 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "advanced.yaml", 10, "", array("fields"=>"field_loop_range"));

$count = sprintf("%d", count($output));
print(">> $count\n");

$arr = explode("|", $output[2]);
$count = sprintf("%d", count($arr));
print(">> $count\n");