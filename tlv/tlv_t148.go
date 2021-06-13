package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T148 struct {
	tlv   *TLV
	bArr  []byte
	j     uint64
	j2    uint64
	j3    uint64
	bArr2 []byte
	bArr3 []byte
}

func NewT148(bArr []byte, j, j2, j3 uint64, bArr2, bArr3 []byte) *T148 {
	return &T148{
		tlv:   NewTLV(0x0148, 0x0000, nil),
		bArr:  bArr,
		j:     j,
		j2:    j2,
		j3:    j3,
		bArr2: bArr2,
		bArr3: bArr3,
	}
}

func (t *T148) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytesN(t.bArr, 0x0020)
	v.EncodeUint32(uint32(t.j))
	v.EncodeUint32(uint32(t.j2))
	v.EncodeUint32(uint32(t.j3))
	v.EncodeBytesN(t.bArr2, 0x0020)
	v.EncodeBytesN(t.bArr3, 0x0020)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T148) Decode(b *bytes.Buffer) error {
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
	j, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.j = uint64(j)
	j2, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.j2 = uint64(j2)
	j3, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.j3 = uint64(j3)
	if t.bArr2, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.bArr3, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
