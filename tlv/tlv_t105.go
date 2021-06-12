package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T105 struct {
	tlv  *TLV
	pic  []byte
	sign []byte
}

func NewT105(pic, sign []byte) *T105 {
	return &T105{
		tlv:  NewTLV(0x0105, 0x0000, nil),
		pic:  pic,
		sign: sign,
	}
}

func (t *T105) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytes(t.sign)
	v.EncodeBytes(t.pic)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T105) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.sign, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.pic, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T105) GetPic() ([]byte, error) {
	return t.pic, nil
}

func (t *T105) GetSign() ([]byte, error) {
	return t.sign, nil
}
