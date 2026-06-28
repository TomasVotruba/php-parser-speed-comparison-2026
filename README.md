# php-parser-comparison

Speed comparison of PHP parsers, run automatically in CI.

Each parser walks the same corpus — a freshly cloned [Laravel framework](https://github.com/laravel/framework) with **all Composer dependencies installed** (`src/` + `vendor/`) — and parses every `.php` file. Each tool runs **10 times** and the **average** wall-clock time is reported.

## Parsers

| Subproject | Parser | Language |
|---|---|---|
| `nikic-PHP-Parser` | [nikic/php-parser](https://github.com/nikic/PHP-Parser) v5 | PHP |
| `ext-ast` | [php-ast](https://github.com/nikic/php-ast) extension | PHP (C ext) |
| `z7zmey-php-parser-dev` | [z7zmey/php-parser](https://github.com/z7zmey/php-parser) | Go |
| `halleck45-go-php-parser` | [halleck45/go-php-parser](https://github.com/Halleck45/go-php-parser) | Go + embedded PHP (cgo) |
| `mago-syntax` | [mago-syntax](https://github.com/carthage-software/mago) v1.42 | Rust |

`mago-syntax` runs in two modes: **parallel** (rayon across all cores, an arena per thread - how [Mago](https://github.com/carthage-software/mago) parses in production) and **single-threaded**. The parallel entry is the only multi-threaded parser here, so for a like-for-like comparison of raw parser speed look at **mago-syntax (single-threaded)**; the parallel number shows real-world throughput.

## Latest results

```
Rank | Parser                  | Avg (10 runs) | vs slowest
   1 | nikic/php-parser (v5)   |      27898 ms |       1.0x
   2 | z7zmey/php-parser       |       4578 ms |       6.1x
   3 | halleck45/go-php-parser |       2721 ms |      10.3x
   4 | ext-ast                 |       2251 ms |      12.4x
```

> Timings come from shared GitHub-hosted runners — good for rough ranking, not precise benchmarking. Live numbers appear in every run's **Summary** page.
>
> **Core count matters.** The `ubuntu-latest` standard runner has only **4 vCPUs** (16 GB RAM). Every parser here runs single-process/single-threaded, so the absolute numbers reflect one core under a noisy-neighbour VM — not bare metal. A machine with more (or faster) cores will post very different timings; only the *relative* ranking is meaningful, and even that can shift with runner contention.

## CI

[`.github/workflows/benchmark.yaml`](.github/workflows/benchmark.yaml) runs on push to `main`, on pull requests, and every 12 hours via cron. One job per parser measures the average; a final `summary` job renders the comparison table into the run summary.
