package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T100 struct {
	tlv              *TLV
	appID            uint64
	subAppID         uint64
	appClientVersion uint32
	mainSigMap       uint32

	ssoVersion uint32
}

func NewT100(appID, subAppID uint64, appClientVersion, mainSigMap, ssoVersion uint32) *T100 {
	return &T100{
		tlv:              NewTLV(0x0100, 0x0000, nil),
		appID:            appID,
		subAppID:         subAppID,
		appClientVersion: appClientVersion,
		mainSigMap:       mainSigMap,

		ssoVersion: ssoVersion,
	}
}

func (t *T100) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0001)
	v.EncodeUint32(t.ssoVersion)
	v.EncodeUint32(uint32(t.appID))
	v.EncodeUint32(uint32(t.subAppID))
	v.EncodeUint32(t.appClientVersion)
	v.EncodeUint32(t.mainSigMap)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T100) Decode(b *bytes.Buffer) error {
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
	if _, err = v.DecodeUint32(); err != nil {
		return err
	}
	appID, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.appID = uint64(appID)
	subAppID, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.subAppID = uint64(subAppID)
	if t.appClientVersion, err = v.DecodeUint32(); err != nil {
		return err
	}
	if t.mainSigMap, err = v.DecodeUint32(); err != nil {
		return err
	}
	return nil
}
