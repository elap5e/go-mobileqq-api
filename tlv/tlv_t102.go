package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T102 struct {
	tlv *TLV
}

func NewT102() *T102 {
	return &T102{
		tlv: NewTLV(0x0102, 0x0000, nil),
	}
}

func (t *T102) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T102) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
