package tlv

import (
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T141 struct {
	tlv   *TLV
	bArr  []byte
	i     uint16
	bArr2 []byte
}

func NewT141(bArr []byte, i uint16, bArr2 []byte) *T141 {
	return &T141{
		tlv:   NewTLV(0x0141, 0x0000, nil),
		bArr:  bArr,
		i:     i,
		bArr2: bArr2,
	}
}

func (t *T141) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0001)
	v.EncodeBytes(t.bArr)
	v.EncodeUint16(t.i)
	v.EncodeBytes(t.bArr2)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T141) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	ver, err := v.DecodeUint16()
	if err != nil {
		return err
	} else if ver != 0x0001 {
		return fmt.Errorf("version 0x%x not support", ver)
	}
	if t.bArr, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.i, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.bArr2, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
