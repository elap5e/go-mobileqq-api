package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T169 struct {
	tlv *TLV
}

func NewT169() *T169 {
	return &T169{
		tlv: NewTLV(0x0169, 0x0000, nil),
	}
}

func (t *T169) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T169) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
