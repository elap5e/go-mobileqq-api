package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T166 struct {
	tlv *TLV
	i   uint8
}

func NewT166(i uint8) *T166 {
	return &T166{
		tlv: NewTLV(0x0166, 0x0000, nil),
		i:   i,
	}
}

func (t *T166) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer([]byte{t.i}))
	t.tlv.Encode(b)
}

func (t *T166) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.i, err = v.DecodeUint8(); err != nil {
		return err
	}
	return nil
}
