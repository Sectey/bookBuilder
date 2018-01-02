package bbld

import (
	"github.com/subchen/go-xmldom"
	"encoding/base64"
)

type fb2binary struct {
	Id          string
	ContentType string
	Binary      []byte
}


func ExtractBinary(doc *xmldom.Document) (res []*fb2binary) {
	for _, n := range doc.Root.Children {
		if n.Name == "binary" {
			b := fb2binary{}
			a := n.GetAttribute("id")
			if a != nil {
				b.Id = a.Value
			}
			a = n.GetAttribute("content-type")
			if a != nil {
				b.ContentType = a.Value
			}
			bt, err := base64.StdEncoding.DecodeString(n.Text)
			if err == nil {
				b.Binary = bt
				res = append(res, &b)
			}
		}
	}
	return res
}