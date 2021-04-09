<?php
include '../vendor/autoload.php';
include '${cls_name}.php';
include 'GPBMetadata/Runtime/Protobuf/${cls_name}.php';

$from = new ${cls_name}();
$from->setName('aaron');

$data = $from->serializeToString();
file_put_contents('data.bin', $data);

$method = $reflect->getMethods();
foreach ($method as $key_m => $value_m) {
    $methodName = $value_m->getName();
    $found = strpos($methodName, 'set');
    if ($found !== false) {
//        var_dump($methodName);

        $comments = $value_m->getDocComment();
//        var_dump($comments);

        // <code>.Address address = 4;</code>
        $pattern = '/<code>\.?(.+?)\s/is';
        preg_match($pattern, $comments, $match);
        echo $match[1] . "\n";
    }
}

// test
$data = file_get_contents('data.bin');
$to = new ${cls_name}();
$to->mergeFromString($data);

echo $to->getName() . PHP_EOL;
