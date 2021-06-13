package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T402 struct {
	tlv *TLV
}

func NewT402() *T402 {
	return &T402{
		tlv: NewTLV(0x0402, 0x0000, nil),
	}
}

func (t *T402) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T402) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
