#!/usr/bin/env php
<?php
/**
[case]
title=reference
cid=0
pid=0

[group]
  1. output >>

[esac]
*/

$output = [];
exec('../build/zd-mac -y definition/refer.yaml -c 3 -field refer -o test/output.txt -f text', $output);
print(">> $output[0]\n");