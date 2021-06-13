package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T161 struct {
	tlv *TLV
}

func NewT161() *T161 {
	return &T161{
		tlv: NewTLV(0x0161, 0x0000, nil),
	}
}

func (t *T161) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T161) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
