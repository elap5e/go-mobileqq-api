package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T118 struct {
	tlv *TLV
}

func NewT118() *T118 {
	return &T118{
		tlv: NewTLV(0x0118, 0x0000, nil),
	}
}

func (t *T118) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T118) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
