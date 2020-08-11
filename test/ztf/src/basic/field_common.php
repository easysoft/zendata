#!/usr/bin/env php
<?php
/**
[case]
title=
cid=1375
pid=7

[group]
 显示10行生成的数据 >>
 验证第12行数据    >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 30, "", array("fields"=>"field_common"));

$count = sprintf("%d", count($output));
print(">> $count\n");
print(">> $output[12]\n");