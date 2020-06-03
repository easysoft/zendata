#!/usr/bin/env php
<?php
/**
[case]
title=field inherit
cid=0
pid=0

[group]
  1. inherit from basic field >>
  2. inherit from reference field >>
[esac]
*/

$output = [];
exec('./zd-mac -d ../test/definition/basic.yaml  -y ../test/definition/inherit.yaml -c 3 -field basic -o ../test/output/output.txt -f text', $output);
print(">> $output[0]\n");

$output = [];
exec('./zd-mac -d ../test/definition/refer.yaml  -y ../test/definition/inherit.yaml -c 3 -field refer2 -o ../test/output/output.txt -f text', $output);
print(">> $output[0]\n");