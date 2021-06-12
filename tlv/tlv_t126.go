package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T126 struct {
	tlv    *TLV
	random []byte
}

func NewT126(random []byte) *T126 {
	return &T126{
		tlv:    NewTLV(0x0126, 0x0000, nil),
		random: random,
	}
}

func (t *T126) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0000)
	v.EncodeBytes(t.random)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T126) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.random, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T126) GetRandom() ([]byte, error) {
	return t.random, nil
}
