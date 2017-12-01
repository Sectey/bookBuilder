package bbld

import (
	"path"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"os"
	"bufio"
	"log"
	"strings"
	"text/template"
	"bytes"
	"path/filepath"
	//"io/ioutil"
)

type Books struct {
	Items []*Book
}

type Book struct {
	BookTitle     string
	Author        Author
	Sequence      Sequence
	ZipFileName      string
	FB2FileName      string
	TxtFileName      string
	CoverpageFile string
	fb2 *Fb2Parser

}

type Author struct {
	FirstName  string
	MiddleName string
	LastName   string
}

type Sequence struct {
	Name string
	Num  string
}

func (me *Book) RenameFb2() error {
	newFileName, err := me.GenFileName(cfg.AudioBook.Fb2FileMask, me.FB2FileName)

	err = os.Rename(me.FB2FileName, newFileName)
	if err != nil {
		return err
	}
	me.FB2FileName = newFileName
	return nil
}

func (me *Book) RenameZip() error {
	if !cfg.AudioBook.ZipRename || cfg.AudioBook.ZipFileMask == "" {
		return nil
	}
	newFileName, err := me.GenFileName(cfg.AudioBook.ZipFileMask, me.ZipFileName)

	err = os.Rename(me.ZipFileName, newFileName)
	if err != nil {
		return err
	}
	me.ZipFileName = newFileName
	return nil
}

func (me *Book) GenFileName(mask string, oldFileName string) (string, error) {
	s := ""
	buf := bytes.NewBufferString(s)
	tmpl, err := template.New("test").Parse(mask)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(buf, *me)

	dir, _ := filepath.Split(oldFileName)
	ext := filepath.Ext(oldFileName)
	newFileName := filepath.Join(dir, buf.String()) + ext
	log.Println("New name:", newFileName)

	return newFileName, err
}

func replaceExt(pathFile string, newExt string) string {
	ext := path.Ext(pathFile)
	return pathFile[0:len(pathFile)-len(ext)] + newExt
}

func (me *Book) SaveTxt() error {
	if me.TxtFileName == "" {
		me.TxtFileName = replaceExt(me.FB2FileName, ".txt")
	}

	text := ""

	nodeList := me.fb2.root.Query("//body")
	for _, c := range nodeList {
		v := c.GetAttributeValue("name")
		if v != "notes" {
			text = text + "\n" + me.fb2.NodeText(c)
		}
	}

	text = me.CorrectText(text)

	//ioutil.WriteFile(me.TxtFileName + ".utf8", []byte(text), 0644)
	writeContent(me.TxtFileName, text)

	return nil
}

func (me *Book) CorrectText(txt string) string {
	t := strings.Replace(txt, "–", "-", -1)
	t = strings.Replace(t, "—", "-", -1)
	t = strings.Replace(t, "–", "-", -1)
	//t = strings.Replace(t, "\xA0", " ", -1)
	t = strings.Replace(t, "\t", " ", -1)
	t = strings.Replace(t, "  ", " ", -1)
	t = strings.Replace(t, "  ", " ", -1)
	t = strings.Replace(t, "  ", " ", -1)
	return t
}

func writeContent(fileName string, str string) (err error) {
	var txtFile *os.File
	if _, err := os.Stat(fileName); os.IsExist(err) {
		os.Remove(fileName)
	}
	txtFile, err = os.Create(fileName)
	if err != nil {
		log.Println(err)
		return err
	}
	w := bufio.NewWriter(txtFile)

	defer txtFile.Close()

	wToWin1251 := transform.NewWriter(w, charmap.Windows1251.NewEncoder())
	wToWin1251.Write ([]byte(str))
	w.Flush()

	return
}