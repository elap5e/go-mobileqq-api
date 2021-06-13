package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T546 struct {
	tlv *TLV
}

func NewT546() *T546 {
	return &T546{
		tlv: NewTLV(0x0546, 0x0000, nil),
	}
}

func (t *T546) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T546) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
