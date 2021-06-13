package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T186 struct {
	tlv  *TLV
	flag bool
}

func NewT186(flag bool) *T186 {
	return &T186{
		tlv:  NewTLV(0x0186, 0x0000, nil),
		flag: flag,
	}
}

func (t *T186) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint8(0x00)
	v.EncodeBool(t.flag)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T186) Decode(b *bytes.Buffer) error {
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
	if t.flag, err = v.DecodeBool(); err != nil {
		return err
	}
	return nil
}

func (t *T186) GetPasswordFlag() (bool, error) {
	return t.flag, nil
}
