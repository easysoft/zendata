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
$output = $zd->create("", "advanced.yaml", 3, "output/advanced.txt");

$zd->decode("advanced.yaml", "output/advanced.txt", "output/advanced.json");

$arr = $zd->readOutput("output/advanced.json");
$content = join($arr, "");
if (strpos($content, "[1|1]") > 0) {
    print(">> found [1|1]\n");
}
if (strpos($content, "[10.0.0.1/'8']") > 0) {
    print(">> found [10.0.0.1/'8']\n");
}