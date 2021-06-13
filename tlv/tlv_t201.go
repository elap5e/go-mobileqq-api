package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T201 struct {
	tlv   *TLV
	bArr  []byte
	bArr2 []byte
	bArr3 []byte
	bArr4 []byte
}

func NewT201(bArr, bArr2, bArr3, bArr4 []byte) *T201 {
	return &T201{
		tlv:   NewTLV(0x0201, 0x0000, nil),
		bArr:  bArr,
		bArr2: bArr2,
		bArr3: bArr3,
		bArr4: bArr4,
	}
}

func (t *T201) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytes(t.bArr)
	v.EncodeBytes(t.bArr2)
	v.EncodeBytes(t.bArr3)
	v.EncodeBytes(t.bArr4)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T201) Decode(b *bytes.Buffer) error {
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
	if t.bArr3, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.bArr4, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
