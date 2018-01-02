package bbld

import (
	"unicode/utf8"
	"golang.org/x/text/transform"
	"errors"
)

type SmartEncoder struct {
}

func (m SmartEncoder) Encode(src []byte, transformer transform.Transformer) (res []byte, err error) {
	enc := s_encoder{src:src, transformer:transformer}

	err = enc.Encode()

	return enc.dst[0:enc.nDst], err
}

type s_encoder struct {
	dst []byte
	src []byte
	nDst int
	nSrc int
	transformer transform.Transformer
}

func (me *s_encoder) next() (res []byte, ok bool) {
	r, size := rune(0), 0

	r = rune(me.src[me.nSrc])

	// Decode a 1-byte rune.
	if r < utf8.RuneSelf {
		size = 1
	} else {
		_, size = utf8.DecodeRune(me.src[me.nSrc:])
	}

	res = me.src[me.nSrc: me.nSrc + size]
	me.nSrc += size
	return res, true
}

func (me *s_encoder) Encode() (err error) {
	me.dst = make([]byte, len(me.src))
	me.transformer.Reset()
	for me.nSrc < len(me.src) {
		s := me.src[me.nSrc:]
		d := me.dst[me.nDst:]
		nd, ns, err1 := me.transformer.Transform(d, s, false)
		me.nSrc += ns
		me.nDst += nd
		if err1 != nil && err1.Error() == "encoding: rune not supported by encoding." {
			if _, ok := me.next(); !ok {
				return errors.New("error")
			}
		} else {
			return err1
		}
	}
	return nil
}
