package tlv

import (
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T185 struct {
	tlv *TLV
	i   uint8
}

func NewT185(i uint8) *T185 {
	return &T185{
		tlv: NewTLV(0x0185, 0x0000, nil),
		i:   i,
	}
}

func (t *T185) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint8(0x01)
	v.EncodeUint8(t.i)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T185) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	ver, err := v.DecodeUint8()
	if err != nil {
		return err
	} else if ver != 0x01 {
		return fmt.Errorf("version 0x%x not support", ver)
	}
	if t.i, err = v.DecodeUint8(); err != nil {
		return err
	}
	return nil
}
