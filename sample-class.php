<?php

namespace App\Service;

use App\Contract\Greeter;

final class GreetingService implements Greeter
{
    private string $prefix;

    public function __construct(string $prefix = 'Hello')
    {
        $this->prefix = $prefix;
    }

    public function greet(string $name): string
    {
        if ($name === '') {
            return $this->prefix . '!';
        }

        return $this->prefix . ', ' . $name . '!';
    }
}
