package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T171 struct {
	tlv *TLV
}

func NewT171() *T171 {
	return &T171{
		tlv: NewTLV(0x0171, 0x0000, nil),
	}
}

func (t *T171) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T171) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
