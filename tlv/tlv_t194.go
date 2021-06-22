package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T194 struct {
	tlv     *TLV
	md5IMSI [16]byte
}

func NewT194(md5IMSI [16]byte) *T194 {
	return &T194{
		tlv:     NewTLV(0x0194, 0x0000, nil),
		md5IMSI: md5IMSI,
	}
}

func (t *T194) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.md5IMSI[:]))
	t.tlv.Encode(b)
}

func (t *T194) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	copy(t.md5IMSI[:], v.Bytes())
	return nil
}
