#!/usr/bin/env php
<?php
/**
[case]
title=repeat
cid=0
pid=0

[group]
  1. output >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '_utils.php';
$cmd = getZDCmd();

$output = [];
exec("$cmd -y ../test/definition/repeat.yaml -c 7 -o ../test/output/output.txt -f text", $output);
print(">> $output[3]\n");