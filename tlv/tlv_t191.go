package tlv

import (
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T191 struct {
	tlv *TLV
	i   uint8
}

func NewT191(i uint8) *T191 {
	return &T191{
		tlv: NewTLV(0x0191, 0x0000, nil),
		i:   i,
	}
}

func (t *T191) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer([]byte{t.i, 0x01}))
	t.tlv.Encode(b)
}

func (t *T191) Decode(b *bytes.Buffer) error {
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
	ver, err := v.DecodeUint8()
	if err != nil {
		return err
	} else if ver != 0x01 {
		return fmt.Errorf("version 0x%x not support", ver)
	}
	return nil
}
