package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T183 struct {
	tlv  *TLV
	salt uint64
}

func NewT183(salt uint64) *T183 {
	return &T183{
		tlv:  NewTLV(0x0183, 0x0000, nil),
		salt: salt,
	}
}

func (t *T183) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint64(t.salt)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T183) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.salt, err = v.DecodeUint64(); err != nil {
		return err
	}
	return nil
}
