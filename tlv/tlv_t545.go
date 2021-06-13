package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T545 struct {
	tlv  *TLV
	bArr [16]byte
}

func NewT545(bArr [16]byte) *T545 {
	return &T545{
		tlv:  NewTLV(0x0545, 0x0000, nil),
		bArr: bArr,
	}
}

func (t *T545) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.bArr[:]))
	t.tlv.Encode(b)
}

func (t *T545) Decode(b *bytes.Buffer) error {
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
