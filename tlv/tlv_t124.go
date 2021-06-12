package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T124 struct {
	tlv   *TLV
	bArr  []byte
	bArr2 []byte
	i     uint16
	bArr3 []byte
	bArr4 []byte
	bArr5 []byte
}

func NewT124(bArr, bArr2 []byte, i uint16, bArr3, bArr4, bArr5 []byte) *T124 {
	return &T124{
		tlv:   NewTLV(0x0124, 0x0000, nil),
		bArr:  bArr,
		bArr2: bArr2,
		i:     i,
		bArr3: bArr3,
		bArr4: bArr4,
		bArr5: bArr5,
	}
}

func (t *T124) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytesN(t.bArr, 0x0010)
	v.EncodeBytesN(t.bArr2, 0x0010)
	v.EncodeBytesN(t.bArr3, 0x0010)
	v.EncodeUint16(t.i)
	v.EncodeBytesN(t.bArr4, 0x0020)
	v.EncodeBytesN(t.bArr5, 0x0010)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T124) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.DecodeUint8(); err != nil {
		return err
	}
	return nil
}
