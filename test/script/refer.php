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
exec('./zd-mac -y ../test/definition/refer.yaml -c 7 -field refer -o ../test/output/output.txt -f text', $output);
print(">> $output[0]\n");