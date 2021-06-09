package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T10D struct {
	tlv *TLV
}

func NewT10D() *T10D {
	return &T10D{
		tlv: NewTLV(0x010d, 0x0000, nil),
	}
}

func (t *T10D) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T10D) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
