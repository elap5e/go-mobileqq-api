package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T403 struct {
	tlv *TLV
}

func NewT403() *T403 {
	return &T403{
		tlv: NewTLV(0x0403, 0x0000, nil),
	}
}

func (t *T403) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T403) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
