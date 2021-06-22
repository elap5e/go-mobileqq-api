package tlv

import (
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T142 struct {
	tlv   *TLV
	apkID []byte
}

func NewT142(apkID []byte) *T142 {
	return &T142{
		tlv:   NewTLV(0x0142, 0x0000, nil),
		apkID: apkID,
	}
}

func (t *T142) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0000)
	v.EncodeBytesN(t.apkID, 0x0020)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T142) Decode(b *bytes.Buffer) error {
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
	if t.apkID, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
