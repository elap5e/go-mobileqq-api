package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T544 struct {
	tlv        *TLV
	uin        uint64
	guid       [16]byte
	sdkVersion string
	typ        uint16
}

func NewT544(uin uint64, guid [16]byte, sdkVersion string, typ uint16) *T544 {
	return &T544{
		tlv:        NewTLV(0x0544, 0x0000, nil),
		uin:        uin,
		guid:       guid,
		sdkVersion: sdkVersion,
		typ:        typ,
	}
}

func (t *T544) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint32(uint32(t.uin))
	v.EncodeBytes(t.guid[:])
	v.EncodeString(t.sdkVersion)
	v.EncodeUint32(uint32(t.typ))
	panic("not implement")
}

func (t *T544) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	panic("not implement")
}
