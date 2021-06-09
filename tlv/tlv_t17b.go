package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T17B struct {
	tlv             *TLV
	availableMsgCnt uint16
	timeLimit       uint16
}

func NewT17B(availableMsgCnt, timeLimit uint16) *T17B {
	return &T17B{
		tlv:             NewTLV(0x017b, 0x0000, nil),
		availableMsgCnt: availableMsgCnt,
		timeLimit:       timeLimit,
	}
}

func (t *T17B) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(t.availableMsgCnt)
	v.EncodeUint16(t.timeLimit)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T17B) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.availableMsgCnt, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.timeLimit, err = v.DecodeUint16(); err != nil {
		return err
	}
	return nil
}

func (t *T17B) GetAvailableMsgCnt() (uint16, error) {
	return t.availableMsgCnt, nil
}

func (t *T17B) GetTimeLimit() (uint16, error) {
	return t.timeLimit, nil
}
