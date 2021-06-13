package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T547 struct {
	tlv  *TLV
	bArr []byte
}

func NewT547(bArr []byte) *T547 {
	return &T547{
		tlv:  NewTLV(0x0547, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T547) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.bArr[:]))
	t.tlv.Encode(b)
}

func (t *T547) Decode(b *bytes.Buffer) error {
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
