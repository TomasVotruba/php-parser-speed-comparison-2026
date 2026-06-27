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

	flag.Parse()

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
