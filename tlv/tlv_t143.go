package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T143 struct {
	tlv *TLV
	d2  []byte
}

func NewT143(d2 []byte) *T143 {
	return &T143{
		tlv: NewTLV(0x0143, 0x0000, nil),
		d2:  d2,
	}
}

func (t *T143) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.d2))
	t.tlv.Encode(b)
}

func (t *T143) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	t.d2 = v.Bytes()
	return nil
}
