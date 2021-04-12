<?php
include '../vendor/autoload.php';
include '${cls_name}.php';
include 'GPBMetadata/Runtime/Protobuf/${cls_name}.php';

$typeArr = array("double","float","int32","int64","uint32","uint64","sint32","sint64",
    "fixed32","fixed64","sfixed32","sfixed64","bool","string","bytes");

$from = new ${cls_name}();
$from->setName('aaron');

$reflect = new ReflectionObject($from);
$method = $reflect->getMethods();
foreach ($method as $key => $value) {
    $methodName = $value->getName();
    $found = strpos($methodName, 'set');
    if ($found !== false) {
        $repeated = false;
        $className = "";

        getProps($value, $repeated, $className);
        printClsField($from, $repeated, $className, $methodName);
    }
}

$data = $from->serializeToString();
file_put_contents('data.bin', $data);

function printClsField($inst, $repeated, $className, $methodName)
{
    global $typeArr;

    $ref = new ReflectionObject($inst);
    echo "object           = " . $ref->getName() . "\n";
    echo "field type       = $className\n";
    echo "field method     = $methodName\n";
    echo "field repeated   = $repeated\n\n";

    if (in_array($className, $typeArr)) {
        return;
    }

    require "./$className.php";

    $propObj = new $className();
    $reflect = new ReflectionObject($propObj);
    $methods = $reflect->getMethods();
    foreach ($methods as $key => $value) {
        $methodName = $value->getName();
        $found = strpos($methodName, 'set');
        if ($found === false) {
            continue;
        }

        $repeated = false;
        $fieldClassName = "";
        getProps($value, $repeated, $fieldClassName);

        if (!in_array($fieldClassName, $typeArr)) {
            printClsField($propObj, $repeated, $fieldClassName, $methodName);
        }
    }
}

function getProps($value, &$repeated, &$className)
{
    $comments = $value->getDocComment();
    // <code>.Address address = 4;</code>
    $pattern = '/<code>(repeated\s)?\.?(.+?)\s/is';
    preg_match($pattern, $comments, $match);
    if (sizeof($match) >= 3) {
        $repeated = $match[1];
        $className = $match[2];
    } else if (sizeof($match) >= 2) {
        $repeated = false;
        $className = $match[1];
    }
    $repeated = trim($repeated);
    if ($repeated === 'repeated')
        $repeated = 1;
    else
        $repeated = 0;

    $className = trim($className);
}

// test
//
//$data = file_get_contents('data.bin');
//$to = new Person();
//$to->mergeFromString($data);
//
//echo $to->getName() . PHP_EOL;
