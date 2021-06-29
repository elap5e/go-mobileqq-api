package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T145 struct {
	tlv  *TLV
	bArr [16]byte
}

func NewT145(bArr [16]byte) *T145 {
	return &T145{
		tlv:  NewTLV(0x0145, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T145) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.bArr[:]))
	t.tlv.Encode(b)
}

func (t *T145) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	copy(t.bArr[:], v.Bytes())
	return nil
}
