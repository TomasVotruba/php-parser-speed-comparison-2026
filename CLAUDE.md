# CLAUDE.md

Benchmark repo comparing PHP parser speed. Each subproject parses the same corpus and reports timing.

## Layout

- `nikic-PHP-Parser/` — PHP, `nikic/php-parser` v5. Composer project. Also pulls `mpdf/mpdf` (only to fatten the corpus, not used by the bench).
- `ext-ast/` — PHP, `php-ast` C extension. Composer requires `ext-ast` (platform), no real packages.
- `z7zmey-php-parser-dev/` — Go, `z7zmey/php-parser` v0.7.2.

The tagged `z7zmey-php-parser/` variant was removed — only the dev one is kept.

## How a benchmark works

- Corpus = the `nikic-PHP-Parser/vendor/` directory. It must exist before any bench runs, so every CI job runs `composer install --working-dir=nikic-PHP-Parser` first.
- Each subproject's `Makefile` has a single `run` target wrapping the parse in `time`.
- `nikic` bench parses its own cwd; `ext-ast` and the Go bench take `../nikic-PHP-Parser` as the path argument.

## Gotchas

- `ext-ast` cannot run without the `ast` PHP extension installed (CI uses `shivammathur/setup-php` with `extensions: ast`). Not installable on a stock box without the extension.
- `ext-ast/bench.php` passes AST version `110` to `ast\parse_file()` — versions below 70 are invalid in php-ast 1.x.
- `nikic/bench.php` uses `(new ParserFactory)->createForNewestSupportedVersion()` — the v4 `create(PREFER_PHP7)` API was removed in v5.
- Go: `z7zmey/php-parser` v0.7.2 changed the API — `php7.NewParser([]byte, version)` and `GetPath()` was removed (the bench prints the file path itself). Older `bytes.Reader`-based code will not compile.
- Built Go binaries (`z7zmey-php-parser-dev/z7zmey-php-parser-dev`) and `vendor/` are gitignored.

## CI

`.github/workflows/benchmark.yaml`: push to `main`, pull requests, cron every 12h.

- One job per parser. Each runs `make run` **10 times**, averages the wall-clock ms, and uploads it as a `duration-*` artifact (`Label|ms` format).
- The `summary` job downloads all artifacts, collects them with `find` (not a glob — files may be nested per artifact), sorts ascending, and renders a fixed-width table (`column -t`) into `$GITHUB_STEP_SUMMARY` (also `tee`'d to the job log).

## Editing the timing table

Keep the artifact line format `Label|ms`. The summary sorts numerically on the second `|`-field, so the label must not contain `|`.
