package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T130 struct {
	tlv    *TLV
	time   []byte
	ipAddr []byte
}

func NewT130(time, ipAddr []byte) *T130 {
	return &T130{
		tlv:    NewTLV(0x0130, 0x0000, nil),
		time:   time,
		ipAddr: ipAddr,
	}
}

func (t *T130) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0000)
	v.EncodeRawBytes(t.time)
	v.EncodeRawBytes(t.ipAddr)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T130) Decode(b *bytes.Buffer) error {
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
	if t.time, err = v.DecodeBytesN(0x0004); err != nil {
		return err
	}
	if t.ipAddr, err = v.DecodeBytesN(0x0004); err != nil {
		return err
	}
	return nil
}

func (t *T130) GetTime() ([]byte, error) {
	return t.time, nil
}

func (t *T130) GetIPAddr() ([]byte, error) {
	return t.ipAddr, nil
}
