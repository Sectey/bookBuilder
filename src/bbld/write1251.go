package bbld

import (
	"os"
	"log"
	"bufio"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/unicode/norm"
	"unicode"
	"strings"
	"bytes"
	"fmt"
)

type Writer1251 struct {
	WriterText
}

func (me *Writer1251) Write(destFileName string, txt string) error  {
	if FileExists(destFileName) {
		os.Remove(destFileName)
	}

	// Create file
	txtFile, err := os.Create(destFileName)
	if err != nil {
		log.Println(err)
		return err
	}
	w := bufio.NewWriter(txtFile)
	defer txtFile.Close()

	writeTxt := txt
	//writeTxt := "́"
	//writeTxt := me.NormalTxt(txt)

	t:= charmap.Windows1251.NewEncoder()
	t.Reset()

	dst, err := SmartEncoder{}.Encode([]byte(writeTxt), t)

	if err != nil {
		fmt.Println(err)
	}

	w.Write(dst)
	w.Flush()
	return err
}

func (me *Writer1251) NormalTxt(txt string) string {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	r := transform.NewReader(strings.NewReader(txt), t)
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(r); err == nil {
		return string(buf.Bytes())
	}
	return ""
}