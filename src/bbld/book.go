package bbld

import (
	"os"
	"path/filepath"
	"strings"
	//"io/ioutil"
	"io/ioutil"
	"log"
	"strconv"
)

type Books struct {
	Items []*Book
}

type Book struct {
	BookTitle     string
	Author        Author
	Sequence      Sequence
	ZipFileName   string
	FB2FileName   string
	TxtFileName   string
	CoverpageFile string
	coverpage     string
	fb2           *Fb2Parser
	WriteBook     bool
	AudioPath     string
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
	newFileName, err := GenFileName(cfg.AudioBook.Fb2FileMask, ExtractDir(me.FB2FileName), me)

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
	newFileName, err := GenFileName(cfg.AudioBook.ZipFileMask, ExtractDir(me.ZipFileName), me)

	err = os.Rename(me.ZipFileName, newFileName)
	if err != nil {
		return err
	}
	me.ZipFileName = newFileName
	return nil
}

func (me *Book) SaveTxt() error {
	//if me.TxtFileName == "" {
	//	me.TxtFileName = replaceExt(me.FB2FileName, ".txt")
	//}
	me.TxtFileName, _ = GenFileName(cfg.AudioBook.TxtFileMask, ExtractDir(me.FB2FileName), me)

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

func (book *Book) SaveCover() {
	num := 0
	num, _ = strconv.Atoi(book.Sequence.Num)
	if cfg.AudioBook.CoverWriteFirstOnly && num > 1 {
		return
	}
	for key, bin := range book.fb2.covers {
		if (key == 1 && book.coverpage == "") || (book.coverpage == bin.Id) {
			filename, err := GenFileName(cfg.AudioBook.CoverFileMask, ExtractDir(book.FB2FileName), book)
			if err != nil {
				log.Println(err)
				continue
			}
			filename = filename + bin.Id
			ioutil.WriteFile(filename, bin.Binary, 0644)
			book.CoverpageFile = filename
		}
	}
}

func (me *Book) CorrectText(txt string) string {
	t := strings.Replace(txt, "–", "-", -1)
	t = strings.Replace(t, "—", "-", -1)
	t = strings.Replace(t, "–", "-", -1)
	//t = strings.Replace(t, "\xA0", " ", -1)
	t = strings.Replace(t, "\t", " ", -1)
	t = strings.Replace(t, "\n", "\r\n", -1)
	t = strings.Replace(t, "  ", " ", -1)
	t = strings.Replace(t, "  ", " ", -1)
	t = strings.Replace(t, "  ", " ", -1)
	return t
}

func writeContent(fileName string, str string) (err error) {
	err = (&Writer1251{}).Write(fileName, str)
	return err
}

func (me *Book) CopyTxtToWritePath() {
	s := me.TxtFileName
	_, fileName := filepath.Split(s)
	newName := filepath.Join(cfg.GetWritingBookPath(), fileName)
	err := os.Rename(me.TxtFileName, newName)
	if (err != nil) {
		me.TxtFileName = newName
	}
}
