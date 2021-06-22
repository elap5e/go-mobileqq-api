package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T16E struct {
	tlv   *TLV
	model []byte
}

func NewT16E(model []byte) *T16E {
	return &T16E{
		tlv:   NewTLV(0x016e, 0x0000, nil),
		model: model,
	}
}

func (t *T16E) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytesN(t.model, 0x0040)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T16E) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	t.model = v.Bytes()
	return nil
}
