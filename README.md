# php-parser-comparison

Speed comparison of PHP parsers, run automatically in CI, every 12 hours.

Each parser walks the same corpus — a freshly cloned [Laravel framework](https://github.com/laravel/framework) with **all Composer dependencies installed** (`src/` + `vendor/`) — and parses every `.php` file. Each tool runs **10 times** and the **average** wall-clock time is reported.

<br>

## Parsers

| Subproject | Parser | Language |
|---|---|---|
| `nikic-PHP-Parser` | [nikic/php-parser](https://github.com/nikic/PHP-Parser) v5 | PHP |
| `ext-ast` | [php-ast](https://github.com/nikic/php-ast) extension | PHP (C ext) |
| `z7zmey-php-parser-dev` | [z7zmey/php-parser](https://github.com/z7zmey/php-parser) | Go |
| `halleck45-go-php-parser` | [halleck45/go-php-parser](https://github.com/Halleck45/go-php-parser) | Go + embedded PHP (cgo) |
| `mago-syntax` | [mago-syntax](https://github.com/carthage-software/mago) v1.42 | Rust |

<br>

## Latest results

Each run produces two tables — every parser pinned to a single core, vs all runner cores available.

### Single core (`taskset -c 0`)

```
Rank | Parser                        | Avg (10 runs) | vs slowest
   1 | nikic/php-parser (v5)         |      29854 ms |       1.0x
   2 | z7zmey/php-parser             |       5941 ms |       5.0x
   3 | halleck45/go-php-parser       |       5048 ms |       5.9x
   4 | ext-ast                       |       2250 ms |      13.3x
   5 | mago-syntax (single-threaded) |       1106 ms |      27.0x
```

<br>

### All cores

```
Rank | Parser                  | Avg (10 runs) | vs slowest
   1 | nikic/php-parser (v5)   |      30092 ms |       1.0x
   2 | z7zmey/php-parser       |       4355 ms |       6.9x
   3 | halleck45/go-php-parser |       2641 ms |      11.4x
   4 | ext-ast                 |       2222 ms |      13.5x
   5 | mago-syntax (parallel)  |        516 ms |      58.3x
```

> Timings come from shared GitHub-hosted runners — good for rough ranking, not precise benchmarking. Live numbers appear in every run's **Summary** page.
>
> **Core count matters.** The `ubuntu-latest` standard runner has only **4 vCPUs** (16 GB RAM). Every parser except `mago-syntax (parallel)` is single-threaded, so its single-core and all-core numbers match — only mago's parallel mode scales with cores. Absolute numbers reflect a noisy-neighbour VM, not bare metal; only the *relative* ranking is meaningful, and even that can shift with runner contention.
