<?php
include '../vendor/autoload.php';
include '${cls_name}.php';
include 'GPBMetadata/Runtime/Protobuf/${cls_name}.php';

$from = new ${cls_name}();
$from->setName('aaron');

$data = $from->serializeToString();
file_put_contents('data.bin', $data);

$reflect = new ReflectionObject($from);
//$props = $reflect->getProperties();
//foreach ($props as $key_p => $value_p) {
//    var_dump($value_p->getName());
//}

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

        $className = $match[1];
        echo "$className\n";

        if (!in_array($className, array("string", "int32"))) {
            require "./$className.php";

            $propObj = new $className();
            $reflect2 = new ReflectionObject($propObj);
            $methods2 = $reflect2->getMethods();
            foreach ($methods2 as $key_m2 => $value_m2) {
                $methodName2 = $value_m2->getName();
                $found2 = strpos($methodName2, 'set');
                if ($found2 !== false) {
                    echo $methodName2 . "\n";
                }
            }
        }
    }
}

// test
$data = file_get_contents('data.bin');
$to = new ${cls_name}();
$to->mergeFromString($data);

echo $to->getName() . PHP_EOL;
