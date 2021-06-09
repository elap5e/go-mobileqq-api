package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T16E struct {
	tlv  *TLV
	bArr []byte
}

func NewT16E(bArr []byte) *T16E {
	return &T16E{
		tlv:  NewTLV(0x016e, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T16E) Encode(b *bytes.Buffer) {
	l := len(t.bArr)
	if l > 0x0040 {
		l = 0x0040
	}
	t.tlv.SetValue(bytes.NewBuffer(t.bArr[:l]))
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
	t.bArr = v.Bytes()
	return nil
}
