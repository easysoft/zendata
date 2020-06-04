#!/usr/bin/env php
<?php

$xml = simplexml_load_file('../test/output/output.xml');
$val = $xml->table->row->col[0];
print(">> $val\n");