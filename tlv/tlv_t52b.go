package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T52B struct {
	tlv    *TLV
	zone   uint16
	mobile string
}

func NewT52B(zone uint16, mobile string) *T52B {
	return &T52B{
		tlv:    NewTLV(0x052b, 0x0000, nil),
		zone:   zone,
		mobile: mobile,
	}
}

func (t *T52B) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint32(0x00000000)
	v.EncodeUint16(t.zone)
	v.EncodeString(t.mobile)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T52B) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.DecodeUint32(); err != nil {
		return err
	}
	if t.zone, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.mobile, err = v.DecodeString(); err != nil {
		return err
	}
	return nil
}

func (t *T52B) GetZone() (uint16, error) {
	return t.zone, nil
}

func (t *T52B) GetMobile() (string, error) {
	return t.mobile, nil
}
