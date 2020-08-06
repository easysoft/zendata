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

$lineArr = $zd->readOutput("output/zt_action.yaml", array(6));
$lines = join("\n", $lineArr);
print(">> $lines\n");