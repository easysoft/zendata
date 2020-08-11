#!/usr/bin/env php
<?php
/**
[case]
title=语言设置
cid=1372
pid=7

[group]
 找到英文文本1 >>
 找到英文文本2 >>
 找到中文文本3 >>
 找到中文文本4 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();

$zd->changeLang("en");

$output = $zd->cmd("-help");
print(">> $output[0]\n");
$output = $zd->cmd("-example");
print(">> $output[0]\n");

$zd->changeLang("zh");

$output = $zd->cmd("-help");
print(">> $output[0]\n");
$output = $zd->cmd("-example");
print(">> $output[0]\n");

$zd->changeLang("en");