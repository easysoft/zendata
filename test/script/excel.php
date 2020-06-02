#!/usr/bin/env php
<?php
/**
[case]
title=query excel
cid=0
pid=0

[group]
  1. usage output >>
  2. basic output >>
  3. step output >>

[esac]
*/

$output = [];
exec('../build/zd-mac -h', $output);
print(">> $output[0]\n");

$output = [];
exec('../build/zd-mac -y definition/basic.yaml -c 15 -field field1 -o test/output.txt -f text', $output);
print(">> $output[0]\n");

$output = [];
exec('../build/zd-mac -y definition/basic.yaml -c 15 -field field2 -o test/output.txt -f text', $output);
print(">> $output[0]\n");