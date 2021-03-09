#!/usr/bin/env php
<?php
/**
[case]
title=查看內置数据定义
cid=1366
pid=7

[group]
 找到预期的文本1 >>
 找到预期的文本2 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();

$output = $zd->cmd("-v");
print(">> $output[0]\n");

$output = $zd->cmd("-view");
print(">> $output[0]\n");