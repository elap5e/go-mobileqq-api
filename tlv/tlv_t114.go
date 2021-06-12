package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T114 struct {
	tlv *TLV
}

func NewT114() *T114 {
	return &T114{
		tlv: NewTLV(0x0114, 0x0000, nil),
	}
}

func (t *T114) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T114) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
