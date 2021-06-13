package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T305 struct {
	tlv *TLV
}

func NewT305() *T305 {
	return &T305{
		tlv: NewTLV(0x0305, 0x0000, nil),
	}
}

func (t *T305) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T305) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
