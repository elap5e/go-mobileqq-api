package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T121 struct {
	tlv *TLV
}

func NewT121() *T121 {
	return &T121{
		tlv: NewTLV(0x0121, 0x0000, nil),
	}
}

func (t *T121) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T121) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
