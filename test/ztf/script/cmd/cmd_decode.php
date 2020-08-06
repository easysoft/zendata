#!/usr/bin/env php
<?php
/**
[case]
title=
cid=0
pid=0

[group]
 >>
 >>

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
