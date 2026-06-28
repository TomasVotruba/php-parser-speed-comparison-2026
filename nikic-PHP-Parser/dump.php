<?php

include_once __DIR__ . "/vendor/autoload.php";

use PhpParser\Error;
use PhpParser\NodeDumper;
use PhpParser\ParserFactory;

$path = $argv[1];
$code = file_get_contents($path);

$parser = (new ParserFactory)->createForNewestSupportedVersion();

try {
    $ast = $parser->parse($code);
} catch (Error $error) {
    echo "Parse error: {$error->getMessage()}\n";
    exit(1);
}

echo (new NodeDumper)->dump($ast) . "\n";
