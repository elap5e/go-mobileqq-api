package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T544 struct {
	tlv *TLV
}

func NewT544() *T544 {
	return &T544{
		tlv: NewTLV(0x0544, 0x0000, nil),
	}
}

func (t *T544) Encode(b *bytes.Buffer) {
	t.tlv.Encode(b)
	panic("not implement")
}

func (t *T544) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	panic("not implement")
}
