package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T107 struct {
	tlv *TLV
	i   uint16
	i2  uint8
	i3  uint16
	i4  uint8
}

func NewT107(i uint16, i2 uint8, i3 uint16, i4 uint8) *T107 {
	return &T107{
		tlv: NewTLV(0x0107, 0x0000, nil),
		i:   i,
		i2:  i2,
		i3:  i3,
		i4:  i4,
	}
}

func (t *T107) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(t.i)
	v.EncodeUint8(t.i2)
	v.EncodeUint16(t.i3)
	v.EncodeUint8(t.i4)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T107) Decode(b *bytes.Buffer) error {
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
	if t.i2, err = v.DecodeUint8(); err != nil {
		return err
	}
	if t.i3, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.i4, err = v.DecodeUint8(); err != nil {
		return err
	}
	return nil
}
