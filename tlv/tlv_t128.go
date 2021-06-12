package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T128 struct {
	tlv   *TLV
	i     uint8
	i2    uint8
	i3    uint8
	i4    uint32
	bArr  []byte
	bArr2 []byte
	bArr3 []byte
}

func NewT128(i, i2, i3 uint8, i4 uint32, bArr, bArr2, bArr3 []byte) *T128 {
	return &T128{
		tlv:   NewTLV(0x0128, 0x0000, nil),
		i:     i,
		i2:    i2,
		i3:    i3,
		i4:    i4,
		bArr:  bArr,
		bArr2: bArr2,
		bArr3: bArr3,
	}
}

func (t *T128) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0000)
	v.EncodeUint8(t.i)
	v.EncodeUint8(t.i2)
	v.EncodeUint8(t.i3)
	v.EncodeUint32(t.i4)
	v.EncodeBytesN(t.bArr, 0x0020)
	v.EncodeBytesN(t.bArr2, 0x0010)
	v.EncodeBytesN(t.bArr3, 0x0010)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T128) Decode(b *bytes.Buffer) error {
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
	if t.i, err = v.DecodeUint8(); err != nil {
		return err
	}
	if t.i2, err = v.DecodeUint8(); err != nil {
		return err
	}
	if t.i3, err = v.DecodeUint8(); err != nil {
		return err
	}
	if t.i4, err = v.DecodeUint32(); err != nil {
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
	return nil
}
