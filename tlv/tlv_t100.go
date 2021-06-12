package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T100 struct {
	tlv   *TLV
	appID uint64
	j2    uint64
	i     uint32
	i2    uint32
}

func NewT100(appID, j2 uint64, i, i2 uint32) *T100 {
	return &T100{
		tlv:   NewTLV(0x0100, 0x0000, nil),
		appID: appID,
		j2:    j2,
		i:     i,
		i2:    i2,
	}
}

func (t *T100) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0001)
	v.EncodeUint32(0x00000011)
	v.EncodeUint32(uint32(t.appID))
	v.EncodeUint32(uint32(t.j2))
	v.EncodeUint32(t.i)
	v.EncodeUint32(t.i2)
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
	j2, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.j2 = uint64(j2)
	if t.i, err = v.DecodeUint32(); err != nil {
		return err
	}
	if t.i2, err = v.DecodeUint32(); err != nil {
		return err
	}
	return nil
}
