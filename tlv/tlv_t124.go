package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T124 struct {
	tlv         *TLV
	osType      []byte
	osVersion   []byte
	networkType uint16
	simOperator []byte
	bArr4       []byte
	apn         []byte
}

func NewT124(osType, osVersion []byte, networkType uint16, simOperator, bArr4, apn []byte) *T124 {
	return &T124{
		tlv:         NewTLV(0x0124, 0x0000, nil),
		osType:      osType,
		osVersion:   osVersion,
		networkType: networkType,
		simOperator: simOperator,
		bArr4:       bArr4,
		apn:         apn,
	}
}

func (t *T124) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytesN(t.osType, 0x0010)
	v.EncodeBytesN(t.osVersion, 0x0010)
	v.EncodeUint16(t.networkType)
	v.EncodeBytesN(t.simOperator, 0x0010)
	v.EncodeBytesN(t.bArr4, 0x0020)
	v.EncodeBytesN(t.apn, 0x0010)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T124) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	_, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}
