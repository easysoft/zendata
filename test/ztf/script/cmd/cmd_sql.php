#!/usr/bin/env php
<?php
/**
[case]
title=从SQL生成配置
cid=1367
pid=7

[group]
 找到预期的id字段 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();

$output = $zd->convertSql("zentao.sql", "output");
$lineArr = $zd->readOutput("output/zt_action.yaml", array(6));
$lines = join("\n", $lineArr);
print(">> $lines\n");