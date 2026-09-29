# php-parser-comparison

Speed comparison of PHP parsers, run automatically in CI every 12 hours.

Each parser walks the same corpus — a freshly cloned [Laravel framework](https://github.com/laravel/framework) with **all Composer dependencies installed** (`src/` + `vendor/`) — and parses every `.php` file. Each tool runs **5 times** and the **average** wall-clock time is reported, along with the **peak memory** (resident set size) of a single run.

<br>

## Parsers

| Subproject | Parser | Language |
|---|---|---|
| `nikic-PHP-Parser` | [nikic/php-parser](https://github.com/nikic/PHP-Parser) v5.7 / v5.8 / v5.9 | PHP |
| `ext-ast` | [php-ast](https://github.com/nikic/php-ast) extension | PHP (C ext) |
| `z7zmey-php-parser-dev` | [z7zmey/php-parser](https://github.com/z7zmey/php-parser) | Go |
| `rector-php-parser-in-go` | [rectorphp/php-parser-in-go](https://github.com/rectorphp/php-parser-in-go) | Go |
| `halleck45-go-php-parser` | [halleck45/go-php-parser](https://github.com/Halleck45/go-php-parser) | Go + embedded PHP (cgo) |
| `mago-syntax` | [mago-syntax](https://github.com/carthage-software/mago) v1.50 | Rust |

<br>

## Latest results

Each run produces two tables — every parser pinned to a single core, vs all runner cores available. `nikic/php-parser` is benchmarked at three versions (v5.7, v5.8, v5.9); they land within run-to-run noise of each other — no version got measurably faster, the library has no perf changes between them.

### Single core (`taskset -c 0`)

```
Rank | Parser                        | Avg (5 runs) | Peak mem | vs slowest
   1 | nikic/php-parser (v5.8)       |     31575 ms |   207 MB |       1.0x
   2 | nikic/php-parser (v5.9)       |     30646 ms |   207 MB |       1.0x
   3 | nikic/php-parser (v5.7)       |     30627 ms |   206 MB |       1.0x
   4 | rectorphp/php-parser-in-go    |      6017 ms |   108 MB |       5.2x
   5 | halleck45/go-php-parser       |      5371 ms |    77 MB |       5.9x
   6 | z7zmey/php-parser             |      3553 ms |    51 MB |       8.9x
   7 | ext-ast                       |      1925 ms |    99 MB |      16.4x
   8 | mago-syntax (single-threaded) |      1193 ms |    71 MB |      26.5x
```

<br>

### All cores

```
Rank | Parser                     | Avg (5 runs) | Peak mem | vs slowest
   1 | nikic/php-parser (v5.9)    |     32051 ms |   207 MB |       1.0x
   2 | nikic/php-parser (v5.8)    |     31027 ms |   207 MB |       1.0x
   3 | nikic/php-parser (v5.7)    |     30337 ms |   207 MB |       1.1x
   4 | rectorphp/php-parser-in-go |      4780 ms |   115 MB |       6.7x
   5 | halleck45/go-php-parser    |      2819 ms |    67 MB |      11.4x
   6 | z7zmey/php-parser          |      2485 ms |    37 MB |      12.9x
   7 | ext-ast                    |      1996 ms |    99 MB |      16.1x
   8 | mago-syntax (parallel)     |       565 ms |   147 MB |      56.7x
```

<br>

## Sample AST dump

Separate CI jobs (`dump-*`) parse one small fixture — [`sample-class.php`](sample-class.php), a ~25-line class — with each parser and print the resulting node tree to that job's **Summary** page. This shows the *shape* of each parser's output (each emits a different format) without wading through the full corpus. Run locally with `make dump` in any subproject.

<br>

Timings come from shared GitHub-hosted runners — good for rough ranking, not precise benchmarking. Live numbers appear in every run's **Summary** page.

**Core count matters.** The `ubuntu-latest` standard runner has only **4 vCPUs** (16 GB RAM). How each parser reacts to extra cores:

- **`mago-syntax (parallel)`** — the only one that actually parses files across cores. Scales **~2.1x** (1193→565 ms) and stays fastest in absolute terms.
- **`nikic`, `ext-ast`** — single-threaded PHP. Single-core and all-core numbers match.
- **`halleck45`, `z7zmey`, `rectorphp/php-parser-in-go`** — parse sequentially, but the Go runtime (GC, scheduler, sysmon) uses extra cores anyway, so pinning to one core (`taskset -c 0`) slows them down. The speedup tracks `GOMAXPROCS`, not the workload — neither does any parallel parsing:
    - `halleck45` gains the most (**~1.9x**: 5371→2819 ms) — Go + cgo around an embedded PHP, so more runtime/allocation work to offload.
    - `z7zmey` is pure Go with less heap churn, so its gain is smaller (**~1.4x**: 3553→2485 ms).
    - `rectorphp/php-parser-in-go` shares the z7zmey lineage and behaves similarly (**~1.25x**: 6017→4780 ms).

**Memory.** The Go parsers are the leanest (`z7zmey` 37–51 MB, `halleck45` 67–77 MB, `rectorphp/php-parser-in-go` 108–115 MB); the PHP tools carry the interpreter's footprint (`nikic` ~207 MB, `ext-ast` ~99 MB). `mago-syntax` is tiny single-threaded (71 MB) but jumps to 147 MB in parallel mode — rayon buffering file contents across worker threads is the memory price of its speed.

Absolute numbers reflect a noisy-neighbour VM, not bare metal; only the *relative* ranking is meaningful, and even that can shift with runner contention.
