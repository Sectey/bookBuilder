package main

import (
	"bbld"
	"path/filepath"

)

const BBLD_FILE_NAME = ".bbld"

func main() {
	path := "c:/Temp/book/Акула пера в Мире Файролла"
	b := bbld.Builder{RootDir: path}
	b.Load(filepath.Join(path, BBLD_FILE_NAME))
	b.Scan()
	b.Save(filepath.Join(path, BBLD_FILE_NAME))
	//
	// tst();
}

