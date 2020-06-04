#!/usr/bin/env php
<?php
if (!function_exists('simplexml_load_file')) {
    $xml = simplexml_load_file('../test/output/output.xml');
    $val = $xml->table->row->col[0];
    print(">> $val\n");
}