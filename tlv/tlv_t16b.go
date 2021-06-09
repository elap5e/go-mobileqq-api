package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T16B struct {
	tlv  *TLV
	list []string
}

func NewT16B(list []string) *T16B {
	return &T16B{
		tlv:  NewTLV(0x016b, 0x0000, nil),
		list: list,
	}
}

func (t *T16B) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(uint16(len(t.list)))
	for _, item := range t.list {
		v.EncodeString(item)
	}
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T16B) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	n, err := v.DecodeUint16()
	if err != nil {
		return err
	}
	t.list = make([]string, n)
	for i := range t.list {
		if t.list[i], err = v.DecodeString(); err != nil {
			return err
		}
	}
	return nil
}
