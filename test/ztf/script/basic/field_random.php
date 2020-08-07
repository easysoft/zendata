#!/usr/bin/env php
<?php
/**
[case]
title=指定随机
cid=1341
pid=7

[group]
 显示10行生成的数据 >>
 验证顺序为随机     >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "default.yaml", 10, "", array("fields"=>"field_random"));

$count = sprintf("%d", count($output));
print(">> $count\n");

$str = join($output, ",");
print("$str\n");

if ($output[0] != 1 || $output[1] != 2 || $output[2] != 3) {
    print(">> not 1,2,3\n");
}