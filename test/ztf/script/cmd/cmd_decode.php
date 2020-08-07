#!/usr/bin/env php
<?php
/**
[case]
title=反向解析数据
cid=1368
pid=7

[group]
 找到特定字符窜1 >>
 找到特定字符窜2 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "test2.yaml", 3, "output/test2.txt");

$zd->decode("test2.yaml", "output/test2.txt", "output/test2.json");

$arr = $zd->readOutput("output/test2.json");
$content = join($arr, "");
if (strpos($content, 'part1_a') > 0) {
    print(">> found part1_a\n");
}
if (strpos($content, 'part3_int_10') > 0) {
    print(">> found part3_int_10\n");
}
