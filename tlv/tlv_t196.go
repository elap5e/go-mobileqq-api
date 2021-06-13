package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T196 struct {
	tlv            *TLV
	bakMobileState uint8
	countryCode    string
	bakMobile      string
}

func NewT196(bakMobileState uint8, countryCode, bakMobile string) *T196 {
	return &T196{
		tlv:            NewTLV(0x0196, 0x0000, nil),
		bakMobileState: bakMobileState,
		countryCode:    countryCode,
		bakMobile:      bakMobile,
	}
}

func (t *T196) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint8(t.bakMobileState)
	v.EncodeString(t.countryCode)
	v.EncodeString(t.bakMobile)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T196) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.bakMobileState, err = v.DecodeUint8(); err != nil {
		return err
	}
	if t.countryCode, err = v.DecodeString(); err != nil {
		return err
	}
	if t.bakMobile, err = v.DecodeString(); err != nil {
		return err
	}
	return nil
}

func (t *T196) GetBakMobileState() (uint8, error) {
	return t.bakMobileState, nil
}

func (t *T196) GetCountryCode() (string, error) {
	return t.countryCode, nil
}

func (t *T196) GetBakMobile() (string, error) {
	return t.bakMobile, nil
}
