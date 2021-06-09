package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T17F struct {
	tlv *TLV
}

func NewT17F() *T17F {
	return &T17F{
		tlv: NewTLV(0x017f, 0x0000, nil),
	}
}

func (t *T17F) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T17F) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
