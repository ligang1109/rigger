<?php
$prjHome = dirname(__DIR__);

$serverConf = json_decode(file_get_contents("$prjHome/conf/server_conf.json"), true);

var_dump($serverConf);
