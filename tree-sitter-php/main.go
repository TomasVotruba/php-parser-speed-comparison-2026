package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_php "github.com/tree-sitter/tree-sitter-php/bindings/go"
)

const usage = "usage: tree-sitter-php-bench <parallel|single> <dir-or-file>"

type parseStats struct {
	parsed     int
	withErrors int
}

type parseResult struct {
	hasError bool
	err      error
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal(usage)
	}

	mode := os.Args[1]
	root := os.Args[2]

	files, err := collectPHPFiles(root)
	if err != nil {
		log.Fatal(err)
	}

	language := sitter.NewLanguage(tree_sitter_php.LanguagePHP())

	var stats parseStats
	switch mode {
	case "parallel":
		stats, err = parseParallel(files, language)
	case "single":
		stats, err = parseSingle(files, language)
	default:
		log.Fatalf("unknown mode %q; %s", mode, usage)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[%s] parsed %d files (%d with parse errors)\n", mode, stats.parsed, stats.withErrors)
}

func collectPHPFiles(root string) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		if filepath.Ext(root) == ".php" {
			return []string{root}, nil
		}
		return nil, nil
	}

	var files []string
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() || filepath.Ext(path) != ".php" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

func parseSingle(files []string, language *sitter.Language) (parseStats, error) {
	parser, err := newParser(language)
	if err != nil {
		return parseStats{}, err
	}
	defer parser.Close()

	var stats parseStats
	for _, path := range files {
		hasError, err := parseOne(parser, path)
		if err != nil {
			return stats, err
		}
		stats.parsed++
		if hasError {
			stats.withErrors++
		}
	}
	return stats, nil
}

func parseParallel(files []string, language *sitter.Language) (parseStats, error) {
	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	if workers > len(files) {
		workers = len(files)
	}
	if workers == 0 {
		return parseStats{}, nil
	}

	parsers := make([]*sitter.Parser, 0, workers)
	for range workers {
		parser, err := newParser(language)
		if err != nil {
			closeParsers(parsers)
			return parseStats{}, err
		}
		parsers = append(parsers, parser)
	}
	defer closeParsers(parsers)

	jobs := make(chan string)
	results := make(chan parseResult, len(files))
	var wg sync.WaitGroup

	for _, parser := range parsers {
		wg.Add(1)
		go func(parser *sitter.Parser) {
			defer wg.Done()
			for path := range jobs {
				hasError, err := parseOne(parser, path)
				results <- parseResult{hasError: hasError, err: err}
			}
		}(parser)
	}

	for _, path := range files {
		jobs <- path
	}
	close(jobs)
	wg.Wait()
	close(results)

	var stats parseStats
	var firstErr error
	for result := range results {
		if result.err != nil && firstErr == nil {
			firstErr = result.err
			continue
		}
		if result.err != nil {
			continue
		}
		stats.parsed++
		if result.hasError {
			stats.withErrors++
		}
	}
	return stats, firstErr
}

func newParser(language *sitter.Language) (*sitter.Parser, error) {
	parser := sitter.NewParser()
	if err := parser.SetLanguage(language); err != nil {
		parser.Close()
		return nil, fmt.Errorf("set PHP tree-sitter language: %w", err)
	}
	return parser, nil
}

func parseOne(parser *sitter.Parser, path string) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	tree := parser.Parse(content, nil)
	if tree == nil {
		return false, fmt.Errorf("tree-sitter returned nil tree for %s", path)
	}
	defer tree.Close()

	return tree.RootNode().HasError(), nil
}

func closeParsers(parsers []*sitter.Parser) {
	for _, parser := range parsers {
		parser.Close()
	}
}
