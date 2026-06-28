<?php

$path = $argv[1];

$ast = ast\parse_file($path, $version = 110);

dumpNode($ast);

function dumpNode($node, string $indent = '')
{
    if ($node instanceof ast\Node) {
        $kind = ast\get_kind_name($node->kind);
        echo $indent . $kind . " (flags={$node->flags})\n";

        foreach ($node->children as $name => $child) {
            echo $indent . "  " . $name . ":\n";
            dumpNode($child, $indent . "    ");
        }

        return;
    }

    echo $indent . var_export($node, true) . "\n";
}
