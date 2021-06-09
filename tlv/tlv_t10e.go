package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T10E struct {
	tlv *TLV
}

func NewT10E() *T10E {
	return &T10E{
		tlv: NewTLV(0x010e, 0x0000, nil),
	}
}

func (t *T10E) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T10E) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
