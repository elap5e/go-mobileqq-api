package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T52C struct {
	tlv *TLV
	i   uint8
	j   uint64
}

func NewT52C(i uint8, j uint64) *T52C {
	return &T52C{
		tlv: NewTLV(0x052c, 0x0000, nil),
		i:   i,
		j:   j,
	}
}

func (t *T52C) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint8(t.i)
	v.EncodeUint64(t.j)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T52C) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.i, err = v.DecodeUint8(); err != nil {
		return err
	}
	if t.j, err = v.DecodeUint64(); err != nil {
		return err
	}
	return nil
}
