package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	PhpParser "github.com/halleck45/go-php-parser/v1"
)

var counter int

func main() {
	defer PhpParser.Shutdown()

	dumpFile := flag.String("dump", "", "parse a single file and dump its AST JSON, then exit")
	flag.Parse()

	if *dumpFile != "" {
		json, ok := PhpParser.ParseFile(*dumpFile, 0, 0)
		if !ok {
			log.Fatalf("failed to parse %s", *dumpFile)
		}
		fmt.Println(json)
		return
	}

	for _, root := range flag.Args() {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Ext(path) != ".php" {
				return nil
			}

			counter++
			PhpParser.ParseFile(path, 0, 0)
			fmt.Printf("[%d] %s\n", counter, path)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
