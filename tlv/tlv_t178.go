package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T178 struct {
	tlv             *TLV
	countryCode     []byte
	mobile          []byte
	smscodeStatus   uint32
	availableMsgCnt uint16
	timeLimit       uint16
}

func NewT178(countryCode, mobile []byte, smscodeStatus uint32, availableMsgCnt, timeLimit uint16) *T178 {
	return &T178{
		tlv:             NewTLV(0x0178, 0x0000, nil),
		countryCode:     countryCode,
		mobile:          mobile,
		smscodeStatus:   smscodeStatus,
		availableMsgCnt: availableMsgCnt,
		timeLimit:       timeLimit,
	}
}

func (t *T178) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytes(t.countryCode)
	v.EncodeBytes(t.mobile)
	v.EncodeUint32(t.smscodeStatus)
	v.EncodeUint16(t.availableMsgCnt)
	v.EncodeUint16(t.timeLimit)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T178) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.countryCode, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.mobile, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.smscodeStatus, err = v.DecodeUint32(); err != nil {
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

func (t *T178) GetCountryCode() ([]byte, error) {
	return t.countryCode, nil
}

func (t *T178) GetMobile() ([]byte, error) {
	return t.mobile, nil
}

func (t *T178) GetSMSCodeStatus() (uint32, error) {
	return t.smscodeStatus, nil
}

func (t *T178) GetAvailableMessageCount() (uint16, error) {
	return t.availableMsgCnt, nil
}

func (t *T178) GetTimeLimit() (uint16, error) {
	return t.timeLimit, nil
}
