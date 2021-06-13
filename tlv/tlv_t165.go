package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T165 struct {
	tlv *TLV
}

func NewT165() *T165 {
	return &T165{
		tlv: NewTLV(0x0165, 0x0000, nil),
	}
}

func (t *T165) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T165) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
