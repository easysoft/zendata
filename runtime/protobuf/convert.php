<?php
include '../vendor/autoload.php';
include '${cls_name}.php';
include 'GPBMetadata/Runtime/Protobuf/${cls_name}.php';

$from = new ${cls_name}();
$from->setName('aaron');

$data = $from->serializeToString();
file_put_contents('data.bin', $data);

// test
$data = file_get_contents('data.bin');
$to = new ${cls_name}();
$to->mergeFromString($data);

echo $to->getName() . PHP_EOL;
