# php-parser-comparison

Speed comparison of PHP parsers, run automatically in CI, every 12 hours.

Each parser walks the same corpus — a freshly cloned [Laravel framework](https://github.com/laravel/framework) with **all Composer dependencies installed** (`src/` + `vendor/`) — and parses every `.php` file. Each tool runs **5 times** and the **average** wall-clock time is reported.

<br>

## Parsers

| Subproject | Parser | Language |
|---|---|---|
| `nikic-PHP-Parser` | [nikic/php-parser](https://github.com/nikic/PHP-Parser) v5 | PHP |
| `ext-ast` | [php-ast](https://github.com/nikic/php-ast) extension | PHP (C ext) |
| `z7zmey-php-parser-dev` | [z7zmey/php-parser](https://github.com/z7zmey/php-parser) | Go |
| `halleck45-go-php-parser` | [halleck45/go-php-parser](https://github.com/Halleck45/go-php-parser) | Go + embedded PHP (cgo) |
| `mago-syntax` | [mago-syntax](https://github.com/carthage-software/mago) v1.42 | Rust |
| `tree-sitter-php` | [tree-sitter-php](https://github.com/tree-sitter/tree-sitter-php) v0.24 | Go binding |

<br>

## Latest results

Each run produces two tables — every parser pinned to a single core, vs all runner cores available.

### Single core (`taskset -c 0`)

```
Rank | Parser                            | Avg (5 runs) | vs slowest
   1 | nikic/php-parser (v5)             |     31407 ms |       1.0x
   2 | tree-sitter-php (single-threaded) |     28546 ms |       1.1x
   3 | z7zmey/php-parser                 |      5666 ms |       5.5x
   4 | halleck45/go-php-parser           |      4481 ms |       7.0x
   5 | ext-ast                           |      2230 ms |      14.1x
   6 | mago-syntax (single-threaded)     |      1025 ms |      30.6x
```

<br>

### All cores

```
Rank | Parser                     | Avg (5 runs) | vs slowest
   1 | nikic/php-parser (v5)      |     30778 ms |       1.0x
   2 | tree-sitter-php (parallel) |     12203 ms |       2.5x
   3 | z7zmey/php-parser          |      4215 ms |       7.3x
   4 | halleck45/go-php-parser    |      2410 ms |      12.8x
   5 | ext-ast                    |      2215 ms |      13.9x
   6 | mago-syntax (parallel)     |       530 ms |      58.1x
```

<br>

Timings come from shared GitHub-hosted runners — good for rough ranking, not precise benchmarking. Live numbers appear in every run's **Summary** page.

**Core count matters.** The `ubuntu-latest` standard runner has only **4 vCPUs** (16 GB RAM). How each parser reacts to extra cores:

- **`mago-syntax (parallel)`** — the only one that parses files across cores. Scales the most.
- **`nikic`, `ext-ast`** — single-threaded PHP. Single-core and all-core numbers match.
- **`halleck45`, `z7zmey`** — parse sequentially, but the Go runtime (GC, scheduler, sysmon) uses extra cores anyway, so pinning to one core (`taskset -c 0`) slows them down. The speedup tracks `GOMAXPROCS`, not the workload — neither does any parallel parsing:
    - `halleck45` gains the most (**~2x**: 5048→2641 ms) — Go + cgo around an embedded PHP, so more runtime/allocation work to offload.
    - `z7zmey` is pure Go with less heap churn, so its gain is smaller (**~1.4x**: 5941→4355 ms).

Absolute numbers reflect a noisy-neighbour VM, not bare metal; only the *relative* ranking is meaningful, and even that can shift with runner contention.
