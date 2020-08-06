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

        if ($default) {
            $cmdStr = str_replace(" -c ", " -d " . $this->workDir . "/" .  $default . " -c ", $cmdStr);
        }

        if (array_key_exists("fields", $options)) {
            $cmdStr .= " -F " .  $options["fields"];
        }
        if (array_key_exists("table", $options)) {
            $cmdStr .= " -table " .  $options["table"];
        }
        if (array_key_exists("trim", $options)) {
            $cmdStr .= " -T ";
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
        $arr = explode("\n", $content);
        if (count($lines) == 0) {
            $ret = array();
            foreach ($arr as $item) {
                $item = trim($item);
                if ($item) {
                    array_push($ret, $item);
                }
            }
            return $ret;
        }

        $ret = array();
        foreach ($lines as $num) {
            array_push($ret, $arr[$num - 1]);
        }

        return $ret;
    }

    public function decode($config, $input, $out)
    {
        $cmdStr = sprintf("%s -D -c %s/%s -i %s/%s -o %s/%s",
            $this->cmdPath, $this->workDir, $config, $this->workDir, $input, $this->workDir, $out);
        print("$cmdStr\n");

        exec($cmdStr, $output);
    }

    public function startService($port, $root)
    {
        $this->stopService($port);

        $cmdStr = sprintf("nohup %s -p %d  >%s/log.txt 2>&1 &", $this->cmdPath, $port, $this->workDir);
        if ($root) {
            $cmdStr = str_replace(" -p ", " -R " . $root . " -p ", $cmdStr);
        }

        print("$cmdStr\n");

        pclose(popen($cmdStr, 'r'));

        exec($cmdStr, $output);

        exec("lsof -i :8848", $output);
        $str = join($output, "\n");
        print("$str\n");
    }
    public function stopService($port)
    {
        $cmdStr = sprintf("kill -9 `lsof -i :%d -t`", $port);
        print("$cmdStr\n");

        exec($cmdStr, $output);
    }
    public function httpGet($port, $default, $conf, $lines, $options = array())
    {
        $url = sprintf("http://127.0.0.1:%d/?d=%s/%s&c=%s/%s&lines=%d",
            $port, $this->workDir, $default, $this->workDir, $conf, $lines);

        if (array_key_exists("root", $options)) {
            $url .= "&root=" . $options["root"];
        }

        print("$url\n");

        $resp = file_get_contents($url);
        print("$resp\n");

        return $resp;
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