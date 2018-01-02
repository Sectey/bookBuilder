package bbld

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Builder struct {
	RootDir string
	Books []*Book
}

func (me *Builder) Add(book *Book) {
	me.Books = append(me.Books, book)
}

func (me *Builder) FindZip(ZipFileName string) *Book {
	for _, value := range me.Books {
		if value.ZipFileName == ZipFileName {
			return value
		}
	}
	r := Book{ZipFileName:ZipFileName}
	me.Add(&r)
	return &r
}

func (me *Builder) Load(fileName string) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("error: %v", err)
		return err
	}

	books := []*Book{}
	err = yaml.Unmarshal([]byte(b), &books)
	if err != nil {
		log.Println("error: %v", err)
	}
	//pbooks := make([]*Book, len(books))
	//for i, value := range books {
	//	pbooks :=
	//}

	return nil
}

func (me *Builder) Save(fileName string) error {
	d, err := yaml.Marshal(&me.Books)
	if err != nil {
		log.Println("error: %v", err)
		return err
	}
	ioutil.WriteFile(fileName, d, 0644)

	return nil
}

func (me *Builder) Scan() error {
	zipFiles := []string{};
	filepath.Walk(me.RootDir, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		if err != nil {
			log.Fatalln(err)
		}
		if info != nil && !info.IsDir() && filepath.Ext(path) == ".zip" {
			zipFiles = append(zipFiles, path)
		}
		return nil
	});
	me.Run(zipFiles)
	return nil
}

func (me *Builder) Run(zipFileNames []string) error {
	for _, fileName := range zipFileNames{
		log.Println(fileName)
		me.PrepareBoor(fileName)
	}
	return nil
}

func (me *Builder) PrepareBoor(ZipFileName string) error {
	book := me.FindZip(ZipFileName)
	if !cfg.AudioBook.RewriteBook && book.WriteBook {
		return nil
	}

	err := me.ExtractFb2(book)
	me.FillBookData(book)
	if err != nil {
		fmt.Println(err)
	}
	book.RenameZip()
	if cfg.AudioBook.Fb2Delete {
		os.Remove(book.FB2FileName)
	}
	book.SaveTxt()
	book.SaveCover()
	book.CopyTxtToWritePath()
	me.FillAudio(book)
	book.WriteBook = true
	return err
}

func (me *Builder) ExtractFb2(book *Book) (err error) {
	i := book.ZipFileName
	book.FB2FileName, err = UnzipOneFile(i, "" )

	return err
}

func (me *Builder) FillBookData(book *Book) error {
	book.fb2 = &Fb2Parser{}
	err := book.fb2.Open(book.FB2FileName)
	if err != nil {
		return err
	}
	err = book.fb2.LoadBookInfo(book)

	book.RenameFb2()
	return err
}

func (me *Builder) FillAudio(book *Book)  {
	if (book.AudioPath != "") {
		return
	}
	audioPath := cfg.GetWritingBookPath()
	_, fl := filepath.Split(book.TxtFileName)
	fl = "!" + strings.Replace(ExtractFilenameWithoutExt(fl), " ", "_", -1)
	audioPath = filepath.Join(audioPath, fl)
	if FileExists(audioPath) {
		book.AudioPath = audioPath
	}
}

func (me *Builder) TrimFb2(book *Book) error {

	return nil
}