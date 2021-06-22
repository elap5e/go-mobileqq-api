package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T188 struct {
	tlv          *TLV
	md5AndroidID [16]byte
}

func NewT188(md5AndroidID [16]byte) *T188 {
	return &T188{
		tlv:          NewTLV(0x0188, 0x0000, nil),
		md5AndroidID: md5AndroidID,
	}
}

func (t *T188) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.md5AndroidID[:]))
	t.tlv.Encode(b)
}

func (t *T188) Decode(b *bytes.Buffer) error {
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
