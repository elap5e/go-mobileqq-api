package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T17E struct {
	tlv *TLV
}

func NewT17E() *T17E {
	return &T17E{
		tlv: NewTLV(0x017e, 0x0000, nil),
	}
}

func (t *T17E) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T17E) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
