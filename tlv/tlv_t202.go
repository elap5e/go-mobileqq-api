package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T202 struct {
	tlv   *TLV
	bArr  []byte
	bArr2 []byte
}

func NewT202(bArr, bArr2 []byte) *T202 {
	return &T202{
		tlv:   NewTLV(0x0202, 0x0000, nil),
		bArr:  bArr,
		bArr2: bArr2,
	}
}

func (t *T202) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	if len(t.bArr) == 0 {
		t.bArr = make([]byte, 16)
	}
	v.EncodeBytesN(t.bArr, 0x0010)
	v.EncodeBytesN(t.bArr2, 0x0020)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T202) Decode(b *bytes.Buffer) error {
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
