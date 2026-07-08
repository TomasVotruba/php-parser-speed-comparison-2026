<?php

include_once __DIR__ . "/vendor/autoload.php";

use PhpParser\Error;
use PhpParser\NodeDumper;
use PhpParser\ParserFactory;
use PhpParser\Parser;

$path = realpath(getcwd() . DIRECTORY_SEPARATOR . $argv[1]);

$iterator = new RecursiveDirectoryIterator($path);
$recursiveIterator = new RecursiveIteratorIterator($iterator);
$parser = (new ParserFactory)->createForNewestSupportedVersion();

/** @var SplFileInfo $path */
foreach ($recursiveIterator as $path) {
    
    if (!$path->isFile()) {
        continue;
    }
    
    if ($path->getExtension() !== 'php') {
        continue;
    }
    
    parse($path, $parser);
    
    echo $path . "\n";
}

function parse(SplFileInfo $path, Parser $parser) {
    try {
        $ast = $parser->parse(file_get_contents($path));
    } catch (Error $error) {
        echo "Parse error: {$error->getMessage()}\n";
    }
}
