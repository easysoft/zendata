#!/usr/bin/env php
<?php
/**
[case]
title=basic
cid=0
pid=0

[group]
  1. usage output >>
  2. basic output >>
  3. step output >>

[esac]
*/

$output = [];
exec('./zd-mac -h', $output);
print(">> $output[0]\n");

$output = [];
exec('./zd-mac -y ../test/definition/basic.yaml -c 3 -field basic -o ../test/output/output.txt -f text', $output);
$count = sprintf("%d", count($output));
print(">> $count\n");
print(">> $output[0]\n");

$output = [];
exec('ls ../test/output/output.txt -f text', $output);
print(">> $output[0]\n");

$output = [];
exec('./zd-mac -y ../test/definition/basic.yaml -c 7 -field step -o ../test/output/output.txt -f text', $output);
print(">> $output[0]\n");