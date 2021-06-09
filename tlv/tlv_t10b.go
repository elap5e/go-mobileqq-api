package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T10B struct {
	tlv *TLV
}

func NewT10B() *T10B {
	return &T10B{
		tlv: NewTLV(0x010b, 0x0000, nil),
	}
}

func (t *T10B) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T10B) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
