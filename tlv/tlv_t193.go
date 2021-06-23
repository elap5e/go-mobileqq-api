package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T193 struct {
	tlv  *TLV
	code []byte
}

func NewT193(code []byte) *T193 {
	return &T193{
		tlv:  NewTLV(0x0193, 0x0000, nil),
		code: code,
	}
}

func (t *T193) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.code))
	t.tlv.Encode(b)
}

func (t *T193) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	t.code = v.Bytes()
	return nil
}
