package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T199 struct {
	tlv      *TLV
	openID   []byte
	payToken []byte
}

func NewT199(openID, payToken []byte) *T199 {
	return &T199{
		tlv:      NewTLV(0x0199, 0x0000, nil),
		openID:   openID,
		payToken: payToken,
	}
}

func (t *T199) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytes(t.openID)
	v.EncodeBytes(t.payToken)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T199) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.openID, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.payToken, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T199) GetOpenID() ([]byte, error) {
	return t.openID, nil
}

func (t *T199) GetPayToken() ([]byte, error) {
	return t.payToken, nil
}
