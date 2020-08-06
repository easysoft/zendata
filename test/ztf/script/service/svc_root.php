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

$port = 8848;
$root = dirname(dirname(dirname(__FILE__)));
$zd->startService($port, $root);

$resp = $zd->httpGet($port, "default.yaml", "test.yaml", 10);

$jsonArr = json_decode($resp,TRUE);

$count = sprintf("%d", count($jsonArr));
print(">> $count\n");

$field = $jsonArr[2]["field_common"];
print(">> $field\n");

$zd->stopService(8848);