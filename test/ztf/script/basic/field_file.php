#!/usr/bin/env php
<?php
/**
[case]
title=文件内容随机
cid=1342
pid=7

[group]
 显示10行生成的数据    >>
 验证数据在预定的数组中 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 10, "", array("fields"=>"field_file"));

$count = sprintf("%d", count($output));
print(">> $count\n");

$str = join($output, ",");
print("$str\n");

$array = array("aaron", "ben", "carl");
if (in_array($output[0], $array)) {
    print(">> $output[0] is in array\n");
}