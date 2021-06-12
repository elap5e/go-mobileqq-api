package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T112 struct {
	tlv  *TLV
	bArr []byte
}

func NewT112(bArr []byte) *T112 {
	return &T112{
		tlv:  NewTLV(0x0112, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T112) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.bArr))
	t.tlv.Encode(b)
}

func (t *T112) Decode(b *bytes.Buffer) error {
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
