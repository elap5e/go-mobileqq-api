package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T174 struct {
	tlv  *TLV
	bArr []byte
}

func NewT174(bArr []byte) *T174 {
	return &T174{
		tlv:  NewTLV(0x0174, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T174) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.bArr))
	t.tlv.Encode(b)
}

func (t *T174) Decode(b *bytes.Buffer) error {
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
