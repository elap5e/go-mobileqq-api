package tlv

import (
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T18 struct {
	tlv   *TLV
	appID uint64
	i1    uint32
	uin   uint64
	i2    uint16
}

func NewT18(appID uint64, i1 uint32, uin uint64, i2 uint16) *T18 {
	return &T18{
		tlv:   NewTLV(0x0018, 0x0000, nil),
		appID: appID,
		i1:    i1,
		uin:   uin,
		i2:    i2,
	}
}

func (t *T18) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0001)
	v.EncodeUint32(0x00000600)
	v.EncodeUint32(uint32(t.appID))
	v.EncodeUint32(t.i1)
	v.EncodeUint32(uint32(t.uin))
	v.EncodeUint16(t.i2)
	v.EncodeUint16(0x0000)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T18) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	pingVer, err := v.DecodeUint16()
	if err != nil {
		return err
	} else if pingVer != 0x0001 {
		return fmt.Errorf("ping version 0x%x not support", pingVer)
	}
	ssoVer, err := v.DecodeUint32()
	if err != nil {
		return err
	} else if ssoVer != 0x00000600 {
		return fmt.Errorf("sso version 0x%x not support", ssoVer)
	}
	appID, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.appID = uint64(appID)
	if t.i1, err = v.DecodeUint32(); err != nil {
		return err
	}
	uin, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.uin = uint64(uin)
	if t.i2, err = v.DecodeUint16(); err != nil {
		return err
	}
	if _, err = v.DecodeUint16(); err != nil {
		return err
	}
	return nil
}
