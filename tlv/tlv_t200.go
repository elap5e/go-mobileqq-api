package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T200 struct {
	tlv   *TLV
	pf    []byte
	pfKey []byte
}

func NewT200(pf, pfKey []byte) *T200 {
	return &T200{
		tlv:   NewTLV(0x0200, 0x0000, nil),
		pf:    pf,
		pfKey: pfKey,
	}
}

func (t *T200) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytes(t.pf)
	v.EncodeBytes(t.pfKey)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T200) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.pf, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.pfKey, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T200) GetPF() ([]byte, error) {
	return t.pf, nil
}

func (t *T200) GetPFKey() ([]byte, error) {
	return t.pfKey, nil
}
