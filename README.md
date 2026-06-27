# php-parser-comparison

Speed comparison of PHP parsers, run automatically in CI.

Each parser walks the same corpus (the `nikic/php-parser` `vendor/` directory) and parses every `.php` file. Each tool runs **10 times** and the **average** wall-clock time is reported.

## Parsers

| Subproject | Parser | Language |
|---|---|---|
| `nikic-PHP-Parser` | [nikic/php-parser](https://github.com/nikic/PHP-Parser) v5 | PHP |
| `ext-ast` | [php-ast](https://github.com/nikic/php-ast) extension | PHP (C ext) |
| `z7zmey-php-parser-dev` | [z7zmey/php-parser](https://github.com/z7zmey/php-parser) | Go |

## Latest results

```
Rank | Parser                | Avg (10 runs)
   1 | ext-ast               |        205 ms
   2 | z7zmey/php-parser     |        267 ms
   3 | nikic/php-parser (v5) |       2237 ms
```

> Timings come from shared GitHub-hosted runners — good for rough ranking, not precise benchmarking. Live numbers appear in every run's **Summary** page.

## Run locally

Each subproject has a `make run` target that wraps the parse in `time`:

```bash
# PHP parsers (need PHP 8.4; ext-ast also needs the `ast` extension)
composer install --working-dir=nikic-PHP-Parser   # provides the corpus
make -C nikic-PHP-Parser run
make -C ext-ast run

# Go parser
go build -o z7zmey-php-parser-dev/z7zmey-php-parser-dev ./z7zmey-php-parser-dev
make -C z7zmey-php-parser-dev run
```

## CI

[`.github/workflows/benchmark.yaml`](.github/workflows/benchmark.yaml) runs on push to `main`, on pull requests, and every 12 hours via cron. One job per parser measures the average; a final `summary` job renders the comparison table into the run summary.
