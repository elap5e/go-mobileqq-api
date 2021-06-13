package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T204 struct {
	tlv *TLV
}

func NewT204() *T204 {
	return &T204{
		tlv: NewTLV(0x0204, 0x0000, nil),
	}
}

func (t *T204) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T204) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
