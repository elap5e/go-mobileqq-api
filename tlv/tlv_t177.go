package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T177 struct {
	tlv *TLV
	j   uint64
	str string
}

func NewT177(j uint64, str string) *T177 {
	return &T177{
		tlv: NewTLV(0x0177, 0x0000, nil),
		j:   j,
		str: str,
	}
}

func (t *T177) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint8(0x01)
	v.EncodeUint32(uint32(t.j))
	v.EncodeString(t.str)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T177) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.DecodeUint8(); err != nil {
		return err
	}
	j, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.j = uint64(j)
	if t.str, err = v.DecodeString(); err != nil {
		return err
	}
	return nil
}
