#!/usr/bin/env php
<?php
/**
[case]
title=field nested
cid=0
pid=0

[group]
  1. output >>
[esac]
*/

$output = [];
exec('./zd-mac -y ../test/definition/nested.yaml -c 3 -field nested', $output);
print(">> $output[0]\n");