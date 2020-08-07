#!/usr/bin/env php
<?php
/**
[case]
title=查看示例
cid=1364
pid=7

[group]
 找到预期的文本1 >>
 找到预期的文本2 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();

$zd->changeLang("en");

$output = $zd->cmd("-e");
print(">> $output[0]\n");
$output = $zd->cmd("-example");
print(">> $output[0]\n");