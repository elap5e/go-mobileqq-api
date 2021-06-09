package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T10A struct {
	tlv  *TLV
	bArr []byte
}

func NewT10A(bArr []byte) *T10A {
	return &T10A{
		tlv:  NewTLV(0x010a, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T10A) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.bArr))
	t.tlv.Encode(b)
}

func (t *T10A) Decode(b *bytes.Buffer) error {
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
