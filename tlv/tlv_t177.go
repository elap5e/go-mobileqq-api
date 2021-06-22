package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T177 struct {
	tlv        *TLV
	buildTime  uint64
	sdkVersion string
}

func NewT177(buildTime uint64, sdkVersion string) *T177 {
	return &T177{
		tlv:        NewTLV(0x0177, 0x0000, nil),
		buildTime:  buildTime,
		sdkVersion: sdkVersion,
	}
}

func (t *T177) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint8(0x01)
	v.EncodeUint32(uint32(t.buildTime))
	v.EncodeString(t.sdkVersion)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T177) Decode(b *bytes.Buffer) error {
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
	buildTime, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.buildTime = uint64(buildTime)
	if t.sdkVersion, err = v.DecodeString(); err != nil {
		return err
	}
	return nil
}
