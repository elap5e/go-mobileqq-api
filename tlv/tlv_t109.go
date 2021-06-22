package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T109 struct {
	tlv          *TLV
	md5AndroidID [16]byte
}

func NewT109(md5AndroidID [16]byte) *T109 {
	return &T109{
		tlv:          NewTLV(0x0109, 0x0000, nil),
		md5AndroidID: md5AndroidID,
	}
}

func (t *T109) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.md5AndroidID[:]))
	t.tlv.Encode(b)
}

func (t *T109) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	copy(t.md5AndroidID[:], v.Bytes())
	return nil
}
