package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T17C struct {
	tlv  *TLV
	bArr []byte
}

func NewT17C(bArr []byte) *T17C {
	return &T17C{
		tlv:  NewTLV(0x017c, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T17C) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytes(t.bArr)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T17C) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.bArr, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
