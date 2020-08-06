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
$output = $zd->create("", "test2.yaml", 10, "output/test2.sql",
    array("fields"=>"test0", "table"=>"tlb_table", "trim"=>"true"));

$lineArr = $zd->readOutput("output/test2.sql");

print(">> $lineArr[2]\n");