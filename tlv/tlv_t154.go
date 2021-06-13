package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T154 struct {
	tlv *TLV
	i   uint32
}

func NewT154(i uint32) *T154 {
	return &T154{
		tlv: NewTLV(0x0154, 0x0000, nil),
		i:   i,
	}
}

func (t *T154) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint32(t.i)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T154) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.i, err = v.DecodeUint32(); err != nil {
		return err
	}
	return nil
}
