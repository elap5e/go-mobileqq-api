package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T173 struct {
	tlv *TLV
}

func NewT173() *T173 {
	return &T173{
		tlv: NewTLV(0x0173, 0x0000, nil),
	}
}

func (t *T173) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T173) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
