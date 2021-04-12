<?php
include '../vendor/autoload.php';
include '${cls_name}.php';
include 'GPBMetadata/Runtime/Protobuf/${cls_name}.php';

$typeArr = array("double","float","int32","int64","uint32","uint64","sint32","sint64",
    "fixed32","fixed64","sfixed32","sfixed64","bool","string","bytes");

$inst = new ${cls_name}();

echo("====================set properties====================\n");
$reflect = new ReflectionObject($inst);
$methods = $reflect->getMethods();
foreach ($methods as $key => $value) {
    $setMethodName = $value->getName();
    $found = strpos($setMethodName, 'set');
    if ($found !== false) {
        $repeated = false;
        $setMethodParamCls = "";

        parserFieldPropsFromComments($value, $repeated, $setMethodParamCls);
        setFieldDefaultValue($inst, $repeated, $setMethodParamCls, $setMethodName);
    }
}

$data = $inst->serializeToString();
file_put_contents('data.bin', $data);

$data = file_get_contents('data.bin');
$decode = new Person();
$decode->mergeFromString($data);

echo("====================print object====================\n");
printObj($decode);

function printObj($obj) {
    $reflect = new ReflectionObject($obj);
    $name = $reflect->getName();
    $methods = $reflect->getMethods();
    foreach ($methods as $key => $value) {
        $methodName = $value->getName();
        $found = strpos($methodName, 'get');
        if ($found !== false && $methodName !== "getClass" && $methodName !== "getType" && $methodName !== "getLegacyClass") {
            $repeated = false;
            $className = "";
            parserFieldPropsFromComments($value, $repeated, $className);
            $var = call_user_func(array($obj, $methodName));

            if (isStandType($className) || gettype($var) === 'integer') {
                printField($name, $methodName, $var, $repeated, $className);
                continue;
            }

            printObj($var);
        }
    }
}

function printField($name, $methodName, $var, $repeated, $className) {
    if ($className === 'bool') {
        if ($var) $var = 'true';
        else $var = 'false';
    }

    if (!$repeated) {
        $name = str_pad($name . '->' . $methodName, 26," ", STR_PAD_RIGHT);
        echo "$name = $var\n";
        return;
    }

    $arr = array();
    foreach ($var as $key => $value) {
        $arr[] = $value;
    }

    $name = str_pad($methodName, 26," ", STR_PAD_RIGHT);
    echo "$name = [" . join(",",$arr) . "] \n";
}

function setFieldDefaultValue(&$inst, $repeated, $setMethodParamCls, $setMethodName)
{
    global $typeArr;

    $ref = new ReflectionObject($inst);
    $objectType = $ref->getName();
    echo "object type      = $objectType\n";
    echo "field type       = $setMethodParamCls\n";
    echo "field method     = $setMethodName\n";
    echo "field repeated   = $repeated\n\n";

    if (isStandType($setMethodParamCls)) { //
        $defaultVal = getDefaultValByType($setMethodParamCls, $repeated);
        call_user_func(array($inst, $setMethodName), $defaultVal);

        return;
    }

    require "./$setMethodParamCls.php";

    $child = new $setMethodParamCls();
    $reflect = new ReflectionObject($child);
    $methods = $reflect->getMethods();
    $isEnum = true;
    foreach ($methods as $key => $value) {
        $setMethodName2 = $value->getName();
        $found = strpos($setMethodName2, 'set');
        if ($found === false) {
            continue;
        }

        $isEnum = false;

        $repeated2 = 0;
        $fieldClassName = "";
        parserFieldPropsFromComments($value, $repeated2, $fieldClassName);
        setFieldDefaultValue($child, $repeated2, $fieldClassName, $setMethodName2);
    }

    if (!$isEnum) {
        $children = $child;
        if ($repeated) {
            $children = array();
            $children[] = $child;
        }

        call_user_func(array($inst, $setMethodName), $children);
    }

}

function parserFieldPropsFromComments($value, &$repeated, &$className)
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

function getDefaultValByType($type, $repeat) {
    if (!$repeat) {
        return getRandValByType($type);
    }

    $count = rand(3, 30);
    $ret = array();
    for ($i = 0; $i < $count; $i++) {
        $item = getRandValByType($type);
        $ret[$i] = $item;
    }

    return $ret;
}

function getRandValByType($type) {
    if ($type === 'bool') { // : string
        $ret = rand(0, 1);
        return $ret;

    } else if ($type === 'string') { // : string
        $r = rand(3, 100);
        $ret = getRandStr($r);
        return $ret;

    } else if ($type === 'float') { // java float: float
        $start = pow(2,7) * -1;
        $end = pow(2,7) - 1;

        $ret = getRandFloat($start, $end);
        return $ret;

    } else if ($type === 'double') { // java double: float
        $start = pow(2,10) * -1;
        $end = pow(2,10) - 1;

        $ret = getRandFloat($start, $end);
        return $ret;

    } else if ($type === 'int32' || $type === 'sint32' || $type === 'sfixed32') { // go int32 : integer
        $start = pow(2,31) * -1;
        $end = pow(2,31) - 1; // 2147483647

        $ret = rand($start, $end);
        return $ret;

    } else if ($type === 'uint32' || $type === 'fixed32') { // go uint32 : integer
        $end = pow(2,32) - 1; // 4294967295

        $ret = rand(0, $end);
        return $ret;

    } else if ($type === 'int64' || $type === 'sint64' || $type === 'sfixed64') { // go int64 : integer/string
        $ret = rand() << 32 | rand(); //    6571882023217245969
        //   -3451450498452162931
        $sign = rand(0, 1);
        if ($sign == 0) {
            $ret *= -1;
        }
        return $ret;

    } else if ($type === 'uint64' || $type === 'fixed64') { // go uint64 : integer
        $ret = rand() << 32 | rand(); //    8252018705439509776
        $ret = ($ret - 1) * 2;
        return $ret;

    }
}

function isStandType($className) {
    global $typeArr;

    return in_array($className, $typeArr);
}

function getRandStr($length = 10) {
    srand(date("s"));
    $chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
    $string = "";
    while(strlen($string) < $length) {
        $string .= substr($chars,(rand()%(strlen($chars))),1);
    }
    return($string);
}

function getRandFloat($min = 0, $max = 1) {
    $rl = mt_rand() / mt_getrandmax();
    return ($min + ($rl * ($max - $min)));
}
