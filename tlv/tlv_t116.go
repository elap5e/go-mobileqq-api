package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T116 struct {
	tlv  *TLV
	i    uint32
	i2   uint32
	jArr []uint64
}

func NewT116(i, i2 uint32, jArr []uint64) *T116 {
	return &T116{
		tlv:  NewTLV(0x0116, 0x0000, nil),
		i:    i,
		i2:   i2,
		jArr: jArr,
	}
}

func (t *T116) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint8(0x00)
	v.EncodeUint32(t.i)
	v.EncodeUint32(t.i2)
	v.EncodeUint8(uint8(len(t.jArr)))
	for i := range t.jArr {
		v.EncodeUint32(uint32(t.jArr[i]))
	}
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T116) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.DecodeUint8(); err != nil {
		return err
	}
	if t.i, err = v.DecodeUint32(); err != nil {
		return err
	}
	if t.i2, err = v.DecodeUint32(); err != nil {
		return err
	}
	l, err := v.DecodeUint8()
	if err != nil {
		return err
	}
	t.jArr = make([]uint64, l)
	for i := range t.jArr {
		j, err := v.DecodeUint32()
		if err != nil {
			return err
		}
		t.jArr[i] = uint64(j)
	}
	return nil
}
