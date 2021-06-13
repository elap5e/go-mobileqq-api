package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T153 struct {
	tlv *TLV
	i   uint16
}

func NewT153(i uint16) *T153 {
	return &T153{
		tlv: NewTLV(0x0153, 0x0000, nil),
		i:   i,
	}
}

func (t *T153) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(t.i)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T153) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.i, err = v.DecodeUint16(); err != nil {
		return err
	}
	return nil
}
