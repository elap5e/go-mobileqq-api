package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T541 struct {
	tlv  *TLV
	flag uint8
}

func NewT541(flag uint8) *T541 {
	return &T541{
		tlv:  NewTLV(0x0541, 0x0000, nil),
		flag: flag,
	}
}

func (t *T541) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer([]byte{t.flag}))
	t.tlv.Encode(b)
}

func (t *T541) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.flag, err = v.DecodeUint8(); err != nil {
		return err
	}
	return nil
}

func (t *T541) GetFlag() (uint8, error) {
	return t.flag, nil
}
