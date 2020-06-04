#!/usr/bin/env php
<?php
/**
[case]
title=generate yaml from database table ddl
cid=0
pid=0

[group]
  1. count >>
  2. line of id field >>
[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '_utils.php';
$cmd = getZDCmd();

$output = [];
exec("$cmd -i ../test/definition/_ddl.sql -o ../test/output", $output);
$str = join("\n", $output);
print("$str\n");

$content = file_get_contents('../test/output/zt_action.yaml');
$arr = explode("\n", $content);
$count = sprintf("%d", count($arr));
print(">> $count\n");

$line = $arr[5];
print(">> $line\n");