package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T508 struct {
	tlv *TLV
}

func NewT508() *T508 {
	return &T508{
		tlv: NewTLV(0x0508, 0x0000, nil),
	}
}

func (t *T508) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
	panic("not implement")
}

func (t *T508) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	panic("not implement")
}
