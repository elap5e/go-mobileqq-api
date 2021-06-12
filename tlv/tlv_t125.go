package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T125 struct {
	tlv     *TLV
	openID  []byte
	openKey []byte
}

func NewT125(openID, openKey []byte) *T125 {
	return &T125{
		tlv:     NewTLV(0x0125, 0x0000, nil),
		openID:  openID,
		openKey: openKey,
	}
}

func (t *T125) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytes(t.openID)
	v.EncodeBytes(t.openKey)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T125) Decode(b *bytes.Buffer) error {
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
	if t.openKey, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T125) GetOpenID() ([]byte, error) {
	return t.openID, nil
}

func (t *T125) GetOpenKey() ([]byte, error) {
	return t.openKey, nil
}
