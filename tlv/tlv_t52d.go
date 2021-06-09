package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T52D struct {
	tlv  *TLV
	bArr []byte
}

func NewT52D(bArr []byte) *T52D {
	return &T52D{
		tlv:  NewTLV(0x052d, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T52D) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.bArr))
	t.tlv.Encode(b)
}

func (t *T52D) Decode(b *bytes.Buffer) error {
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
