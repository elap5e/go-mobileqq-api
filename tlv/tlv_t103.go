package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T103 struct {
	tlv *TLV
}

func NewT103() *T103 {
	return &T103{
		tlv: NewTLV(0x0103, 0x0000, nil),
	}
}

func (t *T103) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
}

func (t *T103) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	return nil
}
