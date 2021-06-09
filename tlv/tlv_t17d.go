package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T17D struct {
	tlv             *TLV
	mbGuideType     uint16
	mbGuideMsg      []byte
	mbGuideInfoType uint16
	mbGuideInfo     []byte
}

func NewT17D(mbGuideType uint16, mbGuideMsg []byte, mbGuideInfoType uint16, mbGuideInfo []byte) *T17D {
	return &T17D{
		tlv:             NewTLV(0x017d, 0x0000, nil),
		mbGuideType:     mbGuideType,
		mbGuideMsg:      mbGuideMsg,
		mbGuideInfoType: mbGuideInfoType,
		mbGuideInfo:     mbGuideInfo,
	}
}

func (t *T17D) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(t.mbGuideType)
	v.EncodeBytes(t.mbGuideMsg)
	v.EncodeUint16(t.mbGuideInfoType)
	v.EncodeBytes(t.mbGuideInfo)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T17D) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.mbGuideType, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.mbGuideMsg, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.mbGuideInfoType, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.mbGuideInfo, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T17D) GetMBGuideType() (uint16, error) {
	return t.mbGuideType, nil
}

func (t *T17D) GetMBGuideMsg() ([]byte, error) {
	return t.mbGuideMsg, nil
}

func (t *T17D) GetMBGuideInfoType() (uint16, error) {
	return t.mbGuideInfoType, nil
}

func (t *T17D) GetMBGuideInfo() ([]byte, error) {
	return t.mbGuideInfo, nil
}
