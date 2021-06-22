package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T116 struct {
	tlv          *TLV
	miscBitmap   uint32
	subSigMap    uint32
	subAppIDList []uint64
}

func NewT116(miscBitmap, subSigMap uint32, subAppIDList []uint64) *T116 {
	return &T116{
		tlv:          NewTLV(0x0116, 0x0000, nil),
		miscBitmap:   miscBitmap,
		subSigMap:    subSigMap,
		subAppIDList: subAppIDList,
	}
}

func (t *T116) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint8(0x00)
	v.EncodeUint32(t.miscBitmap)
	v.EncodeUint32(t.subSigMap)
	v.EncodeUint8(uint8(len(t.subAppIDList)))
	for i := range t.subAppIDList {
		v.EncodeUint32(uint32(t.subAppIDList[i]))
	}
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T116) Decode(b *bytes.Buffer) error {
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
	if t.miscBitmap, err = v.DecodeUint32(); err != nil {
		return err
	}
	if t.subSigMap, err = v.DecodeUint32(); err != nil {
		return err
	}
	l, err := v.DecodeUint8()
	if err != nil {
		return err
	}
	t.subAppIDList = make([]uint64, l)
	for i := range t.subAppIDList {
		j, err := v.DecodeUint32()
		if err != nil {
			return err
		}
		t.subAppIDList[i] = uint64(j)
	}
	return nil
}
