package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T113 struct {
	tlv *TLV
	uin uint64
}

func NewT113(uin uint64) *T113 {
	return &T113{
		tlv: NewTLV(0x0113, 0x0000, nil),
		uin: uin,
	}
}

func (t *T113) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint32(uint32(t.uin))
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T113) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	uin, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.uin = uint64(uin)
	return nil
}

func (t *T113) GetUin() (uint64, error) {
	return t.uin, nil
}
