# CLAUDE.md

Benchmark repo comparing PHP parser speed. Each subproject parses the same corpus and reports timing.

## Layout

- `nikic-PHP-Parser/` — PHP, `nikic/php-parser` v5. Composer project. Also pulls `mpdf/mpdf` (legacy; not used by the bench).
- `ext-ast/` — PHP, `php-ast` C extension. Composer requires `ext-ast` (platform), no real packages.
- `z7zmey-php-parser-dev/` — Go, `z7zmey/php-parser` v0.7.2.
- `halleck45-go-php-parser/` — Go + cgo wrapper around an embedded PHP, `halleck45/go-php-parser`.

The tagged `z7zmey-php-parser/` variant was removed — only the dev one is kept.

## How a benchmark works

- Corpus = a freshly cloned `laravel/framework` **with all Composer dependencies installed** (`git clone --depth 1 ... laravel` then `composer install ... --working-dir=laravel`), parsed as the whole `../laravel` tree (`src/` + `vendor/`). It is gitignored — every CI job clones and installs it.
- After install, the corpus is pruned of intentionally-broken PHP fixtures that hard-crash some parsers (halleck45 exits 255): `rm -rf laravel/tests` and `find laravel/vendor -depth -type d -name tests -exec rm -rf {} +`. Keep these prune steps in sync across CI and the README.
- Each subproject's `Makefile` has a single `run` target wrapping the parse in `time`, pointed at `../laravel`.

## Gotchas

- `ext-ast` cannot run without the `ast` PHP extension installed (CI uses `shivammathur/setup-php` with `extensions: ast`). Not installable on a stock box without the extension.
- `ext-ast/bench.php` passes AST version `110` to `ast\parse_file()` — versions below 70 are invalid in php-ast 1.x.
- `nikic/bench.php` uses `(new ParserFactory)->createForNewestSupportedVersion()` — the v4 `create(PREFER_PHP7)` API was removed in v5.
- `z7zmey/php-parser` v0.7.2 changed the API — `php7.NewParser([]byte, version)` and `GetPath()` was removed (the bench prints the file path itself). Older `bytes.Reader`-based code will not compile.
- Built Go binaries and `vendor/` are gitignored.

## halleck45-go-php-parser build (the tricky one)

cgo wrapper that links an embedded static PHP. The Go module ships **no** prebuilt native libs, and only a **musl** linux release exists (no glibc). Build recipe (mirrored in the CI job):

1. `go.mod` has a relative `replace github.com/halleck45/go-php-parser => ./.halleck-src` — `.halleck-src` is gitignored and must be prepared first:
   - `git clone --depth 1 https://github.com/Halleck45/go-php-parser .halleck-src`
   - download `prebuilt-linux_amd64_musl.tar.gz` from the v0.1.0 release and `tar xzf` it into `.halleck-src/` (lands at `.halleck-src/prebuilt/linux_amd64_musl/`).
2. Build with the musl toolchain (`apt-get install musl-tools`) and the `musl` build tag. The module's cgo flags only add `-Iinclude/php`, so extra absolute `-I` for `main`, `Zend`, `TSRM`, `sapi/embed`, `ext`, `ext/date/lib` must be passed via `CGO_CFLAGS` (absolute paths — cgo compiles in a temp dir, relative `-I` fail):
   ```
   INC="$(pwd)/.halleck-src/prebuilt/linux_amd64_musl/include/php"
   CC=musl-gcc CGO_ENABLED=1 CGO_CFLAGS="-I$INC -I$INC/main -I$INC/Zend -I$INC/TSRM -I$INC/sapi/embed -I$INC/ext -I$INC/ext/date/lib" \
     go build -tags musl -o halleck45-go-php-parser .
   ```
3. At first run the binary fetches runtime libs into `./v1/prebuilt/<target>` (cwd-relative, cached/skipped if present) — gitignored. The timed step does one warm-up run before the loop so the download is not counted.
4. Resulting binary is a musl ELF — needs `/lib/ld-musl-x86_64.so.1` (from `musl-tools`) to run.

## CI

`.github/workflows/benchmark.yaml`: push to `main`, pull requests, cron every 12h.

- One job per parser. Each runs `make run` **10 times**, averages the wall-clock ms, and uploads it as a `duration-*` artifact (`Label|ms` format).
- The `summary` job downloads all artifacts, collects them with `find` (not a glob — files may be nested per artifact), sorts ascending, and renders a fixed-width table (`column -t`) into `$GITHUB_STEP_SUMMARY` (also `tee`'d to the job log).

## Editing the timing table

Keep the artifact line format `Label|ms`. The summary sorts numerically on the second `|`-field, so the label must not contain `|`.
