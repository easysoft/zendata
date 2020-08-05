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

    public function cmd($params)
    {
        $cmdStr = sprintf("%s %s", $this->cmdPath, $params);
        print("$cmdStr\n");

        $output = [];
        exec($cmdStr, $output);
        return $output;
    }

    public function create($default, $conf, $lines, $output, $options = array())
    {
        $cmdStr = sprintf("%s -c %s/%s -n %d -F %s",
            $this->cmdPath, $this->workDir, $conf, $lines, $options["fields"]);
        print("$cmdStr\n");

        $output = [];
        exec($cmdStr, $output);
        return $output;
    }

    public function parse($config, $input)
    {
    }
}