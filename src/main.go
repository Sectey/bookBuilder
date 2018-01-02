package main

import (
	"bbld"
	"path/filepath"

	"flag"
)

const BBLD_FILE_NAME = ".bbld"

func main() {
	dir, _ := filepath.Abs("./")
	d := flag.String("dir", dir, "Directory from book")
	flag.Parse()
	dir = *d
	b := bbld.Builder{RootDir: dir}
	b.Load(filepath.Join(dir, BBLD_FILE_NAME))
	b.Scan()
	b.Save(filepath.Join(dir, BBLD_FILE_NAME))
	//
	// tst();
}

