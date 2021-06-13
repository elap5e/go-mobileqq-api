package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T182 struct {
	tlv       *TLV
	msgCnt    uint16
	timeLimit uint16
}

func NewT182(msgCnt, timeLimit uint16) *T182 {
	return &T182{
		tlv:       NewTLV(0x0182, 0x0000, nil),
		msgCnt:    msgCnt,
		timeLimit: timeLimit,
	}
}

func (t *T182) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint8(0x00)
	v.EncodeUint16(t.msgCnt)
	v.EncodeUint16(t.timeLimit)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T182) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.DecodeUint8(); err != nil {
		return err
	}
	if t.msgCnt, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.timeLimit, err = v.DecodeUint16(); err != nil {
		return err
	}
	return nil
}

func (t *T182) GetMessageCount() (uint16, error) {
	return t.msgCnt, nil
}

func (t *T182) GetTimeLimit() (uint16, error) {
	return t.timeLimit, nil
}
