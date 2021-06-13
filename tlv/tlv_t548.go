package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T548 struct {
	tlv  *TLV
	bArr []byte
}

func NewT548(bArr []byte) *T548 {
	return &T548{
		tlv:  NewTLV(0x0548, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T548) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.bArr[:]))
	t.tlv.Encode(b)
}

func (t *T548) Decode(b *bytes.Buffer) error {
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
