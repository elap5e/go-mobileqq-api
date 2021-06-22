package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T187 struct {
	tlv    *TLV
	md5MAC [16]byte
}

func NewT187(md5MAC [16]byte) *T187 {
	return &T187{
		tlv:    NewTLV(0x0187, 0x0000, nil),
		md5MAC: md5MAC,
	}
}

func (t *T187) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.md5MAC[:]))
	t.tlv.Encode(b)
}

func (t *T187) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	copy(t.md5MAC[:], v.Bytes())
	return nil
}
