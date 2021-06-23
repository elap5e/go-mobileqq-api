package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T191 struct {
	tlv          *TLV
	verifyMethod uint8
}

func NewT191(verifyMethod uint8) *T191 {
	return &T191{
		tlv:          NewTLV(0x0191, 0x0000, nil),
		verifyMethod: verifyMethod,
	}
}

func (t *T191) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer([]byte{t.verifyMethod}))
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
	if t.verifyMethod, err = v.DecodeUint8(); err != nil {
		return err
	}
	return nil
}
