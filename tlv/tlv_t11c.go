package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T11C struct {
	tlv *TLV
}

func NewT11C() *T11C {
	return &T11C{
		tlv: NewTLV(0x011c, 0x0000, nil),
	}
}

func (t *T11C) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T11C) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
