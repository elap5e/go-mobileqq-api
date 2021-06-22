package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T154 struct {
	tlv *TLV
	seq uint16
}

func NewT154(seq uint16) *T154 {
	return &T154{
		tlv: NewTLV(0x0154, 0x0000, nil),
		seq: seq,
	}
}

func (t *T154) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(t.seq)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T154) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.seq, err = v.DecodeUint16(); err != nil {
		return err
	}
	return nil
}
