package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T192 struct {
	tlv *TLV
	url string
}

func NewT192(url string) *T192 {
	return &T192{
		tlv: NewTLV(0x0192, 0x0000, nil),
		url: url,
	}
}

func (t *T192) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer([]byte(t.url)))
	t.tlv.Encode(b)
}

func (t *T192) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	t.url = string(v.Bytes())
	return nil
}

func (t *T192) GetURL() (string, error) {
	return t.url, nil
}
