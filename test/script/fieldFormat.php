#!/usr/bin/env php
<?php
/**
[case]
title=defferent data type with format
cid=0
pid=0

[group]
  1. defferent data type with format >>

[esac]
*/

$output = [];
exec('./zd-mac -y ../test/definition/format.yaml -c 3 -field int,float,char -o ../test/output/output.txt -f text', $output);
print(">> $output[1]\n");