package tlv

import (
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T2 struct {
	tlv  *TLV
	bArr []byte
	sign []byte
}

func NewT2(bArr, sign []byte) *T2 {
	return &T2{
		tlv:  NewTLV(0x0002, 0x0000, nil),
		bArr: bArr,
		sign: sign,
	}
}

func (t *T2) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0000)
	v.EncodeBytes(t.bArr)
	v.EncodeBytes(t.sign)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T2) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	sigVer, err := v.DecodeUint16()
	if err != nil {
		return err
	} else if sigVer != 0x0000 {
		return fmt.Errorf("sig version 0x%x not support", sigVer)
	}
	if t.bArr, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.sign, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T2) GetBArr() ([]byte, error) {
	return t.bArr, nil
}

func (t *T2) GetSign() ([]byte, error) {
	return t.sign, nil
}
