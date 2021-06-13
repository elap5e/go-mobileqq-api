package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T543 struct {
	tlv *TLV
}

func NewT543() *T543 {
	return &T543{
		tlv: NewTLV(0x0543, 0x0000, nil),
	}
}

func (t *T543) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T543) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
