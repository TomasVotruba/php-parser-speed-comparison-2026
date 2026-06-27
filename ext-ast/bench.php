<?php

$path = realpath(getcwd() . DIRECTORY_SEPARATOR . $argv[1]);

$iterator = new RecursiveDirectoryIterator($path);
$recursiveIterator = new RecursiveIteratorIterator($iterator);

/** @var SplFileInfo $path */
foreach ($recursiveIterator as $path) {
    
    if (!$path->isFile()) {
        continue;
    }
    
    if ($path->getExtension() !== 'php') {
        continue;
    }
    
    parse($path);
    
    echo $path . "\n";
}

function parse(SplFileInfo $path) {
    try {
        $ast = ast\parse_file($path, $version=110);
    } catch (Error $error) {
        echo "Parse error: {$error->getMessage()}\n";
    }
}
