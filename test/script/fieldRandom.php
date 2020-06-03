#!/usr/bin/env php
<?php
/**
[case]
title=random
cid=0
pid=0

[group]
  1. output >>

[esac]
*/

$output = [];
exec('./zd-mac -y ../test/definition/basic.yaml -c 7 -field random -o ../test/output/output.txt -f text', $output);
$line = $output[0];
print("Got $line\n");

// 1,2,2,2,8,1
$numbs = explode(",", $line);
$count = sprintf("%d", count($numbs));
print(">> $count\n");

$result = 'not random';
$i = 0;
for($i = 0; $i < $count; $i++) {
  if ($numbs[$i] != $i + 1) { // at lease one not equal sequence numb 1,2,3,4,5,6
    $result = 'is random';
    break;
  }
}
print(">> $result\n");