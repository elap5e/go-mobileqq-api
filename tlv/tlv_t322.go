package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T322 struct {
	tlv *TLV
}

func NewT322() *T322 {
	return &T322{
		tlv: NewTLV(0x0322, 0x0000, nil),
	}
}

func (t *T322) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T322) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
