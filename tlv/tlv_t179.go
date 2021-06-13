package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T179 struct {
	tlv *TLV
	url []byte
}

func NewT179(url []byte) *T179 {
	return &T179{
		tlv: NewTLV(0x0179, 0x0000, nil),
		url: url,
	}
}

func (t *T179) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytes(t.url)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T179) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.url, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T179) GetVerifyURL() ([]byte, error) {
	return t.url, nil
}
