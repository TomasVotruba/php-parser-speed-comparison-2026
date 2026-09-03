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
   1 | nikic/php-parser (v5)         |    25955 ms |  211 MB |       1.0x
   2 | z7zmey/php-parser             |     6153 ms |   49 MB |       4.2x
   3 | rectorphp/php-parser-in-go    |     5548 ms |   92 MB |       4.7x
   4 | halleck45/go-php-parser       |     4505 ms |   76 MB |       5.8x
   5 | ext-ast                       |     2363 ms |   99 MB |      11.0x
   6 | mago-syntax (single-threaded) |      774 ms |   74 MB |      33.5x
```

<br>

### All cores

```
Rank | Parser                        | Avg (5 runs) | Peak mem | vs slowest
   1 | nikic/php-parser (v5)         |    25835 ms |  211 MB |       1.0x
   2 | z7zmey/php-parser             |     4677 ms |   39 MB |       5.5x
   3 | rectorphp/php-parser-in-go    |     4461 ms |   85 MB |       5.8x
   4 | halleck45/go-php-parser       |     2512 ms |   68 MB |      10.3x
   5 | ext-ast                       |     2367 ms |   99 MB |      10.9x
   6 | mago-syntax (parallel)        |      327 ms |  159 MB |      79.0x
```

<br>

## Sample AST dump

Separate CI jobs (`dump-*`) parse one small fixture — [`sample-class.php`](sample-class.php), a ~25-line class — with each parser and print the resulting node tree to that job's **Summary** page. This shows the *shape* of each parser's output (each emits a different format) without wading through the full corpus. Run locally with `make dump` in any subproject.

<br>

Timings come from shared GitHub-hosted runners — good for rough ranking, not precise benchmarking. Live numbers appear in every run's **Summary** page.

**Core count matters.** The `ubuntu-latest` standard runner has only **4 vCPUs** (16 GB RAM). How each parser reacts to extra cores:

- **`mago-syntax (parallel)`** — the only one that actually parses files across cores. Scales **~2.4x** (774→327 ms) and stays fastest in absolute terms.
- **`nikic`, `ext-ast`** — single-threaded PHP. Single-core and all-core numbers match.
- **`halleck45`, `z7zmey`, `rectorphp/php-parser-in-go`** — parse sequentially, but the Go runtime (GC, scheduler, sysmon) uses extra cores anyway, so pinning to one core (`taskset -c 0`) slows them down. The speedup tracks `GOMAXPROCS`, not the workload — neither does any parallel parsing:
    - `halleck45` gains the most (**~1.8x**: 4505→2512 ms) — Go + cgo around an embedded PHP, so more runtime/allocation work to offload.
    - `z7zmey` is pure Go with less heap churn, so its gain is smaller (**~1.3x**: 6153→4677 ms).
    - `rectorphp/php-parser-in-go` shares the z7zmey lineage and behaves similarly (**~1.2x**: 5548→4461 ms).

**Memory.** The Go parsers are the leanest (`z7zmey` 39–49 MB, `halleck45` 68–76 MB, `rectorphp/php-parser-in-go` 85–92 MB); the PHP tools carry the interpreter's footprint (`nikic` ~211 MB, `ext-ast` ~99 MB). `mago-syntax` is tiny single-threaded (74 MB) but jumps to 159 MB in parallel mode — rayon buffering file contents across worker threads is the memory price of its speed.

Absolute numbers reflect a noisy-neighbour VM, not bare metal; only the *relative* ranking is meaningful, and even that can shift with runner contention.
