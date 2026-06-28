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
| `php-parser-in-go` | [TomasVotruba/php-parser-in-go](https://github.com/TomasVotruba/php-parser-in-go) (private) | Go |

<br>

## Latest results

Each run produces two tables — every parser pinned to a single core, vs all runner cores available.

### Single core (`taskset -c 0`)

```
   # | Parser                            | Avg (5 runs) | vs slowest
   1 | nikic/php-parser (v5)             |     29369 ms |       1.0x
   2 | tree-sitter-php (single-threaded) |     28660 ms |       1.0x
   3 | z7zmey/php-parser                 |      5917 ms |       5.0x
   4 | php-parser-in-go                  |      5697 ms |       5.2x
   5 | halleck45/go-php-parser           |      4178 ms |       7.0x
   6 | ext-ast                           |      2191 ms |      13.4x
   7 | mago-syntax (single-threaded)     |      1109 ms |      26.5x
```

<br>

### All cores

```
   # | Parser                     | Avg (5 runs) | vs slowest
   1 | nikic/php-parser (v5)      |     29192 ms |       1.0x
   2 | tree-sitter-php (parallel) |     12349 ms |       2.4x
   3 | z7zmey/php-parser          |      4478 ms |       6.5x
   4 | php-parser-in-go           |      3549 ms |       8.2x
   5 | halleck45/go-php-parser    |      2221 ms |      13.1x
   6 | ext-ast                    |      2200 ms |      13.3x
   7 | mago-syntax (parallel)     |       524 ms |      55.7x
```

<br>

Timings come from shared GitHub-hosted runners — good for rough ranking, not precise benchmarking. Live numbers appear in every run's **Summary** page.

**Core count matters.** The `ubuntu-latest` standard runner has only **4 vCPUs** (16 GB RAM). How each parser reacts to extra cores:

- **`mago-syntax (parallel)`, `tree-sitter-php (parallel)`, `php-parser-in-go`** — the three that actually parse files across cores. `tree-sitter-php` scales the most (**~2.3x**: 28660→12349 ms); `php-parser-in-go` scales **~1.6x** (5697→3549 ms); `mago-syntax` scales **~2.1x** (1109→524 ms) and stays fastest in absolute terms.
- **`nikic`, `ext-ast`** — single-threaded PHP. Single-core and all-core numbers match.
- **`halleck45`, `z7zmey`** — parse sequentially, but the Go runtime (GC, scheduler, sysmon) uses extra cores anyway, so pinning to one core (`taskset -c 0`) slows them down. The speedup tracks `GOMAXPROCS`, not the workload — neither does any parallel parsing:
    - `halleck45` gains the most (**~1.9x**: 4178→2221 ms) — Go + cgo around an embedded PHP, so more runtime/allocation work to offload.
    - `z7zmey` is pure Go with less heap churn, so its gain is smaller (**~1.3x**: 5917→4478 ms).

Absolute numbers reflect a noisy-neighbour VM, not bare metal; only the *relative* ranking is meaningful, and even that can shift with runner contention.
