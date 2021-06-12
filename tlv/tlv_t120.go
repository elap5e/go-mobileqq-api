package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T120 struct {
	tlv *TLV
}

func NewT120() *T120 {
	return &T120{
		tlv: NewTLV(0x0120, 0x0000, nil),
	}
}

func (t *T120) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T120) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
