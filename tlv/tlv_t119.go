package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T119 struct {
	tlv *TLV
}

func NewT119() *T119 {
	return &T119{
		tlv: NewTLV(0x0119, 0x0000, nil),
	}
}

func (t *T119) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T119) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
