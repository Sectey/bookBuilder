package bbld

import (
	"github.com/subchen/go-xmldom"
	"log"
	"errors"
	"launchpad.net/xmlpath"
	"fmt"
	"os"
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
)

type Fb2Parser struct {
	doc *xmldom.Document
	root *xmldom.Node
	xpath *xmlpath.Node
	covers []*fb2binary
}

func (me *Fb2Parser) Open(filename string) (err error) {
	if (filename == "") {
		return errors.New("FB2 empty")
	}

	// удаляем ссылку
	me.DelLinkA(filename)


	me.doc, err = xmldom.ParseFile(filename)
	if err != nil {
		log.Printf("error %v", err)
		return err
	}
	me.root = me.doc.Root

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	me.xpath, err = xmlpath.Parse(file)
	if err != nil {
		log.Fatal(err)
	}


	//name := root.GetAttributeValue("name")
	//time := root.GetAttributeValue("time")
	//fmt.Printf("testsuite: name=%v, time=%v\n", name, time)
	//
	//for _, node := range root.GetChildren("testcase") {
	//	name := node.GetAttributeValue("name")
	//	time := node.GetAttributeValue("time")
	//	fmt.Printf("testcase: name=%v, time=%v\n", name, time)
	//}
	return nil
}

func (me *Fb2Parser) LoadBookInfo(book *Book) (err error) {
	if me.root == nil {
		err = errors.New("Not load fb2 book")
		return err
	}

	if book.BookTitle == "" {
		book.BookTitle, _ = elementValue(me.xpath, "FictionBook/description/title-info/book-title")
		if cfg.Fb2Parsing.BookTitleIsLastWords {
			words := strings.Split(book.BookTitle, ".")
			book.BookTitle = strings.TrimSpace(words[len(words)-1])
		}
	}

	if book.Author.FirstName == "" {
		book.Author.FirstName, _ = elementValue(me.xpath, "FictionBook/description/title-info/author/first-name")
	}
	if book.Author.MiddleName == "" {
		book.Author.MiddleName, _ = elementValue(me.xpath, "FictionBook/description/title-info/author/middle-name")
	}
	if book.Author.LastName == "" {
		book.Author.LastName, _ = elementValue(me.xpath, "FictionBook/description/title-info/author/last-name")
	}

	if book.Sequence.Name == "" {
		book.Sequence.Name, _ = elementValue(me.xpath, "FictionBook/description/title-info/sequence/@name")
	}
	if book.Sequence.Num == "" {
		book.Sequence.Num, _ = elementValue(me.xpath, "FictionBook/description/title-info/sequence/@number")
	}

	book.coverpage, _ = elementValue(me.xpath, "FictionBook/description/title-info/coverpage/image/@href")
	if book.coverpage != "" {
		book.coverpage = book.coverpage[1:]
	}

	me.covers = ExtractBinary(me.doc)
	return nil
}

func _nodeText(buffer *bytes.Buffer, node *xmldom.Node, path string) {
	// до текта
	switch node.Name {
	case "v":
		buffer.WriteString(node.Text)
		break
	case "emphasis":
		buffer.WriteString(" ")
		buffer.WriteString(node.Text)
		break
	case "p":
		buffer.WriteString(" ")
		buffer.WriteString(node.Text)
		break
	case "title":
		buffer.WriteString(" \n")
		break
	case "epigraph":
		buffer.WriteString("- Эпиграф -\n")
		break
	case "empty-line":
		buffer.WriteString(" \n \n")
	default:
	}

	for _, value := range node.Children {
		_nodeText(buffer, value, path + "/" + node.Name)
	}

	// после
	switch node.Name {
	case "text-author":
		buffer.WriteString("\n")
		break
	case "stanza":
		buffer.WriteString("\n")
		break
	case "<stanza>":
		buffer.WriteString("\n")
		break
	case "v":
		buffer.WriteString("\n")
		break
	case "emphasis":
		buffer.WriteString("\n")
		break
	case "p":
		buffer.WriteString("\n")
		break
	case "title":
		buffer.WriteString(node.Text)
		buffer.WriteString("\n")
		break
	default:
		buffer.WriteString(node.Text)
	}
}

func (me *Fb2Parser) NodeText(node *xmldom.Node) string {
	var buffer bytes.Buffer

	_nodeText(&buffer, node, "")

	return buffer.String()
}

func elementValue(root *xmlpath.Node, path string) (value string, ok bool) {
	if v, ok := xmlpath.MustCompile(path).String(root); ok {
		return v, ok
	} else {
		fmt.Println("not found:", path)
		return "", ok
	}
}

func (me *Fb2Parser) DelLinkA(filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	s := string(b[:])
	//s = "\\w"
	r, _ := regexp.Compile("<a[ :\\w\\d\\s=\\\"#>\\[\\]</.]+/a>")
	s = r.ReplaceAllString(s, "")
	ioutil.WriteFile(filename, []byte(s), 0644)
}