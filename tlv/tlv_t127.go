package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T127 struct {
	tlv   *TLV
	bArr  []byte
	bArr2 []byte
}

func NewT127(bArr, bArr2 []byte) *T127 {
	return &T127{
		tlv:   NewTLV(0x0127, 0x0000, nil),
		bArr:  bArr,
		bArr2: bArr2,
	}
}

func (t *T127) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0000)
	v.EncodeBytes(t.bArr)
	v.EncodeBytes(t.bArr2)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T127) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.DecodeUint16(); err != nil {
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
