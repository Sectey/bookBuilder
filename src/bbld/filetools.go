package bbld

import (
	"bytes"
	"path/filepath"
	"text/template"
	"os"
)

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ExtractDir(filename string) string {
	dir, _ := filepath.Split(filename)
	return dir
}

func ExtractFilename(filename string) string {
	_, file := filepath.Split(filename)
	return file
}

func ExtractFilenameWithoutExt(filename string) string {
	fl := ExtractFilename(filename)
	ext := filepath.Ext(fl)
	return fl[0:len(fl)-len(ext)]
}

func ReplaceExt(pathFile string, newExt string) string {
	ext := filepath.Ext(pathFile)
	return pathFile[0:len(pathFile)-len(ext)] + newExt
}

func GenFileName(mask string, dir string, obj interface{}) (string, error) {
	s := ""
	buf := bytes.NewBufferString(s)
	tmpl, err := template.New("test").Parse(mask)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(buf, obj)

	newFileName := filepath.Join(dir, buf.String())

	return newFileName, err
}
