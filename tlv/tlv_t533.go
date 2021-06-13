package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T533 struct {
	tlv   *TLV
	bArr  []byte
	bArr2 []byte
}

func NewT533(bArr, bArr2 []byte) *T533 {
	return &T533{
		tlv:   NewTLV(0x0533, 0x0000, nil),
		bArr:  bArr,
		bArr2: bArr2,
	}
}

func (t *T533) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytes(t.bArr)
	v.EncodeBytes(t.bArr2)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T533) Decode(b *bytes.Buffer) error {
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
	if t.bArr2, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
