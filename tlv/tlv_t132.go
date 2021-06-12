package tlv

import (
	"math/rand"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T132 struct {
	tlv         *TLV
	accessToken []byte
	openID      []byte
}

func NewT132(accessToken, openID []byte) *T132 {
	return &T132{
		tlv:         NewTLV(0x0132, 0x0000, nil),
		accessToken: accessToken,
		openID:      openID,
	}
}

func (t *T132) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeRawBytes(t.accessToken)
	v.EncodeUint32(rand.Uint32())
	v.EncodeRawBytes(t.openID)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T132) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.accessToken, err = v.DecodeBytes(); err != nil {
		return err
	}
	if _, err = v.DecodeUint32(); err != nil {
		return err
	}
	if t.openID, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T132) GetAccessToken() ([]byte, error) {
	return t.accessToken, nil
}

func (t *T132) GetOpenID() ([]byte, error) {
	return t.openID, nil
}
