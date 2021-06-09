package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T8 struct {
	tlv *TLV
	i1  uint16
	i2  uint32
	i3  uint16
}

func NewT8(i1 uint16, i2 uint32, i3 uint16) *T8 {
	return &T8{
		tlv: NewTLV(0x0008, 0x0000, nil),
		i1:  i1,
		i2:  i2,
		i3:  i3,
	}
}

func (t *T8) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(t.i1)
	v.EncodeUint32(t.i2)
	v.EncodeUint16(t.i3)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T8) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.i1, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.i2, err = v.DecodeUint32(); err != nil {
		return err
	}
	if t.i3, err = v.DecodeUint16(); err != nil {
		return err
	}
	return nil
}
