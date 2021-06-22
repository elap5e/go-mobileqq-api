package tlv

import (
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T141 struct {
	tlv         *TLV
	simOperator []byte
	networkType uint16
	apn         []byte
}

func NewT141(simOperator []byte, networkType uint16, apn []byte) *T141 {
	return &T141{
		tlv:         NewTLV(0x0141, 0x0000, nil),
		simOperator: simOperator,
		networkType: networkType,
		apn:         apn,
	}
}

func (t *T141) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0001)
	v.EncodeBytes(t.simOperator)
	v.EncodeUint16(t.networkType)
	v.EncodeBytes(t.apn)
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
	if t.simOperator, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.networkType, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.apn, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
