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

## Latest results

```
Rank | Parser                  | Avg (10 runs) | vs slowest
   1 | nikic/php-parser (v5)   |       1718 ms |       1.0x
   2 | z7zmey/php-parser       |        261 ms |       6.6x
   3 | halleck45/go-php-parser |        252 ms |       6.8x
   4 | ext-ast                 |        208 ms |       8.3x
```

> Timings come from shared GitHub-hosted runners — good for rough ranking, not precise benchmarking. Live numbers appear in every run's **Summary** page.

## Run locally

First get the corpus (once, at the repo root) — clone Laravel and install its dependencies:

```bash
git clone --depth 1 https://github.com/laravel/framework laravel
composer install --no-interaction --ignore-platform-reqs --working-dir=laravel
# drop intentionally-broken PHP fixtures that crash some parsers
rm -rf laravel/tests
find laravel/vendor -depth -type d -name tests -exec rm -rf {} +
```

Then run each subproject's `make run` target (wraps the parse in `time`):

```bash
# PHP parsers (need PHP 8.4; ext-ast also needs the `ast` extension)
composer install --working-dir=nikic-PHP-Parser
make -C nikic-PHP-Parser run
make -C ext-ast run

# z7zmey (Go)
go build -o z7zmey-php-parser-dev/z7zmey-php-parser-dev ./z7zmey-php-parser-dev
make -C z7zmey-php-parser-dev run
```

`halleck45-go-php-parser` is a cgo wrapper around an embedded PHP and needs extra setup (musl toolchain + prebuilt native libs). See the `halleck45-go-php-parser` job in the workflow and [CLAUDE.md](CLAUDE.md) for the exact build steps.

## CI

[`.github/workflows/benchmark.yaml`](.github/workflows/benchmark.yaml) runs on push to `main`, on pull requests, and every 12 hours via cron. One job per parser measures the average; a final `summary` job renders the comparison table into the run summary.
