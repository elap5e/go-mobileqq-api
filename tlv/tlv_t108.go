package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T108 struct {
	tlv  *TLV
	ksid []byte
}

func NewT108(ksid []byte) *T108 {
	return &T108{
		tlv:  NewTLV(0x0108, 0x0000, nil),
		ksid: ksid,
	}
}

func (t *T108) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.ksid))
	t.tlv.Encode(b)
}

func (t *T108) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	t.ksid = v.Bytes()
	return nil
}
