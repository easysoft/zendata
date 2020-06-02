#!/usr/bin/env php
<?php
/**
[case]
title=query excel
cid=0
pid=0

[group]
  1. output >>

[esac]
*/

$output = [];
exec('./zd-mac -y ../test/definition/refer.yaml -c 7 -field excel -o ../test/output/output.txt -f text', $output);
print(">> $output[0]\n");