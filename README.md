# php-parser-comparison

Speed comparison of PHP parsers, run automatically in CI every 12 hours.

Each parser walks the same corpus — a freshly cloned [Laravel framework](https://github.com/laravel/framework) with **all Composer dependencies installed** (`src/` + `vendor/`) — and parses every `.php` file. Each tool runs **5 times** and the **average** wall-clock time is reported, along with the **peak memory** (resident set size) of a single run.

<br>

## Parsers

| Subproject | Parser | Language |
|---|---|---|
| `nikic-PHP-Parser` | [nikic/php-parser](https://github.com/nikic/PHP-Parser) v5 | PHP |
| `ext-ast` | [php-ast](https://github.com/nikic/php-ast) extension | PHP (C ext) |
| `z7zmey-php-parser-dev` | [z7zmey/php-parser](https://github.com/z7zmey/php-parser) | Go |
| `rector-php-parser-in-go` | [rectorphp/php-parser-in-go](https://github.com/rectorphp/php-parser-in-go) | Go |
| `halleck45-go-php-parser` | [halleck45/go-php-parser](https://github.com/Halleck45/go-php-parser) | Go + embedded PHP (cgo) |
| `mago-syntax` | [mago-syntax](https://github.com/carthage-software/mago) v1.42 | Rust |

<br>

## Latest results

Each run produces two tables — every parser pinned to a single core, vs all runner cores available.

### Single core (`taskset -c 0`)

```
Rank | Parser                        | Avg (5 runs) | Peak mem | vs slowest
   1 | nikic/php-parser (v5)         |     28346 ms |   204 MB |       1.0x
   2 | z7zmey/php-parser             |      6054 ms |    29 MB |       4.7x
   3 | halleck45/go-php-parser       |      4660 ms |    75 MB |       6.1x
   4 | ext-ast                       |      2532 ms |    99 MB |      11.2x
   5 | mago-syntax (single-threaded) |      1177 ms |    85 MB |      24.1x
```

<br>

### All cores

```
Rank | Parser                  | Avg (5 runs) | Peak mem | vs slowest
   1 | nikic/php-parser (v5)   |     28861 ms |   204 MB |       1.0x
   2 | z7zmey/php-parser       |      4501 ms |    32 MB |       6.4x
   3 | ext-ast                 |      2503 ms |    99 MB |      11.5x
   4 | halleck45/go-php-parser |      2483 ms |    65 MB |      11.6x
   5 | mago-syntax (parallel)  |       553 ms |   153 MB |      52.2x
```

<br>

## Sample AST dump

Separate CI jobs (`dump-*`) parse one small fixture — [`sample-class.php`](sample-class.php), a ~25-line class — with each parser and print the resulting node tree to that job's **Summary** page. This shows the *shape* of each parser's output (each emits a different format) without wading through the full corpus. Run locally with `make dump` in any subproject.

<br>

Timings come from shared GitHub-hosted runners — good for rough ranking, not precise benchmarking. Live numbers appear in every run's **Summary** page.

**Core count matters.** The `ubuntu-latest` standard runner has only **4 vCPUs** (16 GB RAM). How each parser reacts to extra cores:

- **`mago-syntax (parallel)`** — the only one that actually parses files across cores. Scales **~2.1x** (1177→553 ms) and stays fastest in absolute terms.
- **`nikic`, `ext-ast`** — single-threaded PHP. Single-core and all-core numbers match.
- **`halleck45`, `z7zmey`** — parse sequentially, but the Go runtime (GC, scheduler, sysmon) uses extra cores anyway, so pinning to one core (`taskset -c 0`) slows them down. The speedup tracks `GOMAXPROCS`, not the workload — neither does any parallel parsing:
    - `halleck45` gains the most (**~1.9x**: 4660→2483 ms) — Go + cgo around an embedded PHP, so more runtime/allocation work to offload.
    - `z7zmey` is pure Go with less heap churn, so its gain is smaller (**~1.3x**: 6054→4501 ms).

**Memory.** The Go parsers are the leanest (`z7zmey` 29–32 MB, `halleck45` 65–75 MB); the PHP tools carry the interpreter's footprint (`nikic` ~204 MB, `ext-ast` ~99 MB). `mago-syntax` is tiny single-threaded (85 MB) but jumps to 153 MB in parallel mode — rayon buffering file contents across worker threads is the memory price of its speed.

Absolute numbers reflect a noisy-neighbour VM, not bare metal; only the *relative* ranking is meaningful, and even that can shift with runner contention.
