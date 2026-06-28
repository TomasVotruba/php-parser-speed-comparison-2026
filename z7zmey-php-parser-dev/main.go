package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/karrick/godirwalk"
	"github.com/yookoala/realpath"
	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/visitor"
)

var counter int

func main() {
	dumpFile := flag.String("dump", "", "parse a single file and dump its AST, then exit")
	flag.Parse()

	if *dumpFile != "" {
		dumpAST(*dumpFile)
		return
	}

	processPath(flag.Args())
}

func dumpAST(path string) {
	content, err := ioutil.ReadFile(path)
	checkErr(err)

	p := php7.NewParser(content, "7.4")
	p.Parse()

	dumper := &visitor.Dumper{Writer: os.Stdout, Indent: "  "}
	p.GetRootNode().Walk(dumper)
}

func processPath(pathList []string) {
	files := []string{}

	for _, path := range pathList {
		real, err := realpath.Realpath(path)
		checkErr(err)

		s, err := os.Stat(real)
		checkErr(err)

		if !s.IsDir() {
			files = append(files, real)
			continue
		}

		godirwalk.Walk(real, &godirwalk.Options{
			Unsorted: true,
			Callback: func(osPathname string, de *godirwalk.Dirent) error {
				if !de.IsDir() && filepath.Ext(osPathname) == ".php" {
					files = append(files, osPathname)
				}
				return nil
			},
			ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
				return godirwalk.SkipNode
			},
		})
	}

	for _, p := range files {
		parseFile(p)
	}
}

func parseFile(path string) {
	counter++

	content, err := ioutil.ReadFile(path)
	checkErr(err)

	p := php7.NewParser(content, "7.4")
	p.Parse()

	fmt.Println("[" + strconv.Itoa(counter) + "] " + path)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
