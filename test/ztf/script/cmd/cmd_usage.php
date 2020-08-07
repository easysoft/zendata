#!/usr/bin/env php
<?php
/**
[case]
title=查看帮助
cid=1363
pid=7

[group]
 找到预期的文本1 >>
 找到预期的文本2 >>
 找到预期的文本3 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();

$zd->changeLang("en");

$output = $zd->cmd("");
print(">> $output[0]\n");
$output = $zd->cmd("-h");
print(">> $output[0]\n");
$output = $zd->cmd("-help");
print(">> $output[0]\n");