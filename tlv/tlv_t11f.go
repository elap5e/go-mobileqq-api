package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T11F struct {
	tlv     *TLV
	chgTime int64
	tkPri   uint32
}

func NewT11F(chgTime int64, tkPri uint32) *T11F {
	return &T11F{
		tlv:     NewTLV(0x011f, 0x0000, nil),
		chgTime: chgTime,
		tkPri:   tkPri,
	}
}

func (t *T11F) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint32(uint32(t.chgTime))
	v.EncodeUint32(t.tkPri)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T11F) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	chgTime, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.chgTime = int64(chgTime)
	if t.tkPri, err = v.DecodeUint32(); err != nil {
		return err
	}
	return nil
}

func (t *T11F) GetAppID() (int64, error) {
	return t.chgTime, nil
}

func (t *T11F) GetTKPri() (uint32, error) {
	return t.tkPri, nil
}
