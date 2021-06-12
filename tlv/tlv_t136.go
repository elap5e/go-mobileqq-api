package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T136 struct {
	tlv *TLV
}

func NewT136() *T136 {
	return &T136{
		tlv: NewTLV(0x0136, 0x0000, nil),
	}
}

func (t *T136) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T136) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
