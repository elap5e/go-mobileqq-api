package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T108 struct {
	tlv  *TLV
	bArr []byte
}

func NewT108(bArr []byte) *T108 {
	return &T108{
		tlv:  NewTLV(0x0108, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T108) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.bArr))
	t.tlv.Encode(b)
}

func (t *T108) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	t.bArr = v.Bytes()
	return nil
}
