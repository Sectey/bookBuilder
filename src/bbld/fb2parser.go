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
)

type Fb2Parser struct {
	doc *xmldom.Document
	root *xmldom.Node
	xpath *xmlpath.Node
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
	// find all children
	//fmt.Printf("children = %v\n", len(node.Query("//*")))

	book.BookTitle, _ = elementValue(me.xpath, "FictionBook/description/title-info/book-title")

	book.Author.FirstName, _ = elementValue(me.xpath, "FictionBook/description/title-info/author/first-name")
	book.Author.MiddleName, _ = elementValue(me.xpath, "FictionBook/description/title-info/author/middle-name")
	book.Author.LastName, _ = elementValue(me.xpath, "FictionBook/description/title-info/author/last-name")

	book.Sequence.Name, _ = elementValue(me.xpath, "FictionBook/description/title-info/sequence/@name")
	book.Sequence.Num, _ = elementValue(me.xpath, "FictionBook/description/title-info/sequence/@number")
	return nil
}

func _nodeText(buffer *bytes.Buffer, node *xmldom.Node, path string) {
	// до текта
	switch node.Name {
	case "v":
		buffer.WriteString(node.Text)
		break
	case "p":
		buffer.WriteString(" ")
		buffer.WriteString(node.Text)
		break
	case "title":
		buffer.WriteString("\n")
		break
	case "epigraph":
		buffer.WriteString("- Эпиграф -\n")
		break
	case "empty-line":
		buffer.WriteString("\n\n")
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
	r, _ := regexp.Compile("<a[ :\\w\\d\\s=\"#>\\[\\]<]+/a>")
	s = r.ReplaceAllString(s, "")
	ioutil.WriteFile(filename, []byte(s), 0644)
}