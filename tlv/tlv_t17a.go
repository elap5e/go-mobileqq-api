package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T17A struct {
	tlv *TLV
	j   uint64
}

func NewT17A(j uint64) *T17A {
	return &T17A{
		tlv: NewTLV(0x017a, 0x0000, nil),
		j:   j,
	}
}

func (t *T17A) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint32(uint32(t.j))
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T17A) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	j, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.j = uint64(j)
	return nil
}
