#!/usr/bin/env php
<?php
/**
[case]
title=列出內置数据定义
cid=1365
pid=7

[group]
 找到预期文本1 >>
 找到预期文本2 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();

$output = $zd->cmd("-l");
print(">> $output[0]\n");

$output = $zd->cmd("-list");
print(">> $output[0]\n");