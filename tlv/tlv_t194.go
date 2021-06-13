package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T194 struct {
	tlv  *TLV
	bArr []byte
}

func NewT194(bArr []byte) *T194 {
	return &T194{
		tlv:  NewTLV(0x0194, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T194) Encode(b *bytes.Buffer) {
	if len(t.bArr) == 0 {
		t.bArr = make([]byte, 16)
	}
	t.tlv.SetValue(bytes.NewBuffer(t.bArr))
	t.tlv.Encode(b)
}

func (t *T194) Decode(b *bytes.Buffer) error {
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
