<?php

include_once __DIR__ . DIRECTORY_SEPARATOR . 'utils.php';

class zendata
{
    var $cmdPath;
    var $workDir;

    public function __construct()
    {
        global $config;

        $this->cmdPath = getZdCmd();
        $this->workDir = $config['zd']['workDir'];
    }

    public function create($default, $conf, $lines, $output, $options = array())
    {
        $cmdStr = sprintf("%s -c %s/%s -n %d",
            $this->cmdPath, $this->workDir, $conf, $lines);
        if (count($options) > 0 && $options["fields"]) {
            $cmdStr .= " -F " .  $options["fields"];
        }
        if ($output) {
            $cmdStr .= " -o " . $this->workDir . "/" . $output;
        }
        print("$cmdStr\n");

        $output = [];
        exec($cmdStr, $output);
        return $output;
    }

    public function convertSql($file, $dir, $options = array())
    {
        $cmdStr = sprintf("%s -i %s/%s -o %s/%s", $this->cmdPath, $this->workDir, $file, $this->workDir, $dir);
        print("$cmdStr\n");

        $output = [];
        exec($cmdStr, $output);
        return $output;
    }

    public function readOutput($file, $lines=array())
    {
        $filePath = sprintf("%s/%s", $this->workDir, $file);
        print("$filePath\n");

        $content = file_get_contents($filePath);
        if (count($lines) == 0) {
            return $content;
        }

        $ret = array();
        $arr = explode("\n", $content);
        foreach ($lines as $num) {
            array_push($ret, $arr[$num - 1]);
        }

        return $ret;
    }

    public function decode($config, $input, $out)
    {
        $cmdStr = sprintf("-D -c %s/%s -i %s/%s -o %s/%s",
            $this->cmdPath, $this->workDir, $config, $this->workDir, $input, $this->workDir, $out);
        print("$cmdStr\n");

        exec($cmdStr, $output);
    }

    public function cmd($params)
    {
        $cmdStr = sprintf("%s %s", $this->cmdPath, $params);
        print("$cmdStr\n");

        $output = [];
        exec($cmdStr, $output);
        return $output;
    }

    public function changeLang($lang)
    {
        $filename = '';

        $filename = sprintf("conf/zdata.conf");
        $content = file_get_contents($filename);
        $fp = fopen($filename, "w");

        $content = str_replace("en", $lang, $content);
        $content = str_replace("zh", $lang, $content);

        fwrite($fp, $content);
        fclose($fp);
    }
}