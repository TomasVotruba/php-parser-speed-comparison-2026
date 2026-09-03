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
	"github.com/rectorphp/php-parser-in-go/pkg/conf"
	"github.com/rectorphp/php-parser-in-go/pkg/parser"
	"github.com/rectorphp/php-parser-in-go/pkg/version"
	"github.com/rectorphp/php-parser-in-go/pkg/visitor/dumper"
	"github.com/yookoala/realpath"
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

func phpVersion() *version.Version {
	v, err := version.New("8.3")
	checkErr(err)
	return v
}

func dumpAST(path string) {
	content, err := ioutil.ReadFile(path)
	checkErr(err)

	root, err := parser.Parse(content, conf.Config{Version: phpVersion()})
	checkErr(err)

	dumper.NewDumper(os.Stdout).Dump(root)
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

	v := phpVersion()
	for _, p := range files {
		parseFile(p, v)
	}
}

func parseFile(path string, v *version.Version) {
	counter++

	content, err := ioutil.ReadFile(path)
	checkErr(err)

	parser.Parse(content, conf.Config{Version: v})

	fmt.Println("[" + strconv.Itoa(counter) + "] " + path)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
