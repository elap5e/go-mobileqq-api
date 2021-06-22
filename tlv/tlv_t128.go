package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T128 struct {
	tlv           *TLV
	isGuidFileNil bool
	isGuidGenSucc bool
	isGuidChanged bool
	guidFlag      uint32
	model         []byte
	guid          []byte
	brand         []byte
}

func NewT128(isGuidFileNil, isGuidGenSucc, isGuidChanged bool, guidFlag uint32, model, guid, brand []byte) *T128 {
	return &T128{
		tlv:           NewTLV(0x0128, 0x0000, nil),
		isGuidFileNil: isGuidFileNil,
		isGuidGenSucc: isGuidGenSucc,
		isGuidChanged: isGuidChanged,
		guidFlag:      guidFlag,
		model:         model,
		guid:          guid,
		brand:         brand,
	}
}

func (t *T128) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0000)
	v.EncodeBool(t.isGuidFileNil)
	v.EncodeBool(t.isGuidGenSucc)
	v.EncodeBool(t.isGuidChanged)
	v.EncodeUint32(t.guidFlag)
	v.EncodeBytesN(t.model, 0x0020)
	v.EncodeBytesN(t.guid, 0x0010)
	v.EncodeBytesN(t.brand, 0x0010)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T128) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.isGuidFileNil, err = v.DecodeBool(); err != nil {
		return err
	}
	if t.isGuidGenSucc, err = v.DecodeBool(); err != nil {
		return err
	}
	if t.isGuidChanged, err = v.DecodeBool(); err != nil {
		return err
	}
	if t.guidFlag, err = v.DecodeUint32(); err != nil {
		return err
	}
	if t.model, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.guid, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.brand, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
