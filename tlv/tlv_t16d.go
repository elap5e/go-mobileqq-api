package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T16D struct {
	tlv *TLV
}

func NewT16D() *T16D {
	return &T16D{
		tlv: NewTLV(0x016d, 0x0000, nil),
	}
}

func (t *T16D) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T16D) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
