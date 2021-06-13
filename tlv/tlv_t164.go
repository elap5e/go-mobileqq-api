package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T164 struct {
	tlv *TLV
}

func NewT164() *T164 {
	return &T164{
		tlv: NewTLV(0x0164, 0x0000, nil),
	}
}

func (t *T164) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T164) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
