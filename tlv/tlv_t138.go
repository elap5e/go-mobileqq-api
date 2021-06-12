package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T138 struct {
	tlv          *TLV
	a2ChgTime    uint32
	lsKeyChgTime uint32
	sKeyChgTime  uint32
	vKeyChgTime  uint32
	a8ChgTime    uint32
	stWebChgTime uint32
	d2ChgTime    uint32
	sidChgTime   uint32
}

func NewT138(a2ChgTime, lsKeyChgTime, sKeyChgTime, vKeyChgTime, a8ChgTime, stWebChgTime, d2ChgTime, sidChgTime uint32) *T138 {
	return &T138{
		tlv:          NewTLV(0x0138, 0x0000, nil),
		a2ChgTime:    a2ChgTime,
		lsKeyChgTime: lsKeyChgTime,
		sKeyChgTime:  sKeyChgTime,
		vKeyChgTime:  vKeyChgTime,
		a8ChgTime:    a8ChgTime,
		stWebChgTime: stWebChgTime,
		d2ChgTime:    d2ChgTime,
		sidChgTime:   sidChgTime,
	}
}

func (t *T138) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint32(0x00000008)
	v.EncodeUint16(0x010a)
	v.EncodeUint32(t.a2ChgTime)
	v.EncodeUint32(0x00000000)
	v.EncodeUint16(0x011c)
	v.EncodeUint32(t.lsKeyChgTime)
	v.EncodeUint32(0x00000000)
	v.EncodeUint16(0x0120)
	v.EncodeUint32(t.sKeyChgTime)
	v.EncodeUint32(0x00000000)
	v.EncodeUint16(0x0136)
	v.EncodeUint32(t.vKeyChgTime)
	v.EncodeUint32(0x00000000)
	v.EncodeUint16(0x0102)
	v.EncodeUint32(t.a8ChgTime)
	v.EncodeUint32(0x00000000)
	v.EncodeUint16(0x0103)
	v.EncodeUint32(t.stWebChgTime)
	v.EncodeUint32(0x00000000)
	v.EncodeUint16(0x0143)
	v.EncodeUint32(t.d2ChgTime)
	v.EncodeUint32(0x00000000)
	v.EncodeUint16(0x0164)
	v.EncodeUint32(t.sidChgTime)
	v.EncodeUint32(0x00000000)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T138) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	l, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	for i := 0; i < int(l); i++ {
		tt, err := v.DecodeUint16()
		if err != nil {
			return err
		}
		time, err := v.DecodeUint32()
		if err != nil {
			return err
		}
		if _, err = v.DecodeUint32(); err != nil {
			return err
		}
		switch tt {
		case 0x010a:
			t.a2ChgTime = time
		case 0x011c:
			t.lsKeyChgTime = time
		case 0x0120:
			t.sKeyChgTime = time
		case 0x0136:
			t.vKeyChgTime = time
		case 0x0102:
			t.a8ChgTime = time
		case 0x0103:
			t.stWebChgTime = time
		case 0x0143:
			t.d2ChgTime = time
		case 0x0164:
			t.sidChgTime = time
		}
	}
	return nil
}

func (t *T138) GetA2ChgTime() (uint32, error) {
	return t.a2ChgTime, nil
}

func (t *T138) GetLSKeyChgTime() (uint32, error) {
	return t.lsKeyChgTime, nil
}

func (t *T138) GetSKeyChgTime() (uint32, error) {
	return t.sKeyChgTime, nil
}

func (t *T138) GetVKeyChgTime() (uint32, error) {
	return t.vKeyChgTime, nil
}

func (t *T138) GetA8ChgTime() (uint32, error) {
	return t.a8ChgTime, nil
}

func (t *T138) GetSTWebChgTime() (uint32, error) {
	return t.stWebChgTime, nil
}

func (t *T138) GetD2ChgTime() (uint32, error) {
	return t.d2ChgTime, nil
}

func (t *T138) GetSidChgTime() (uint32, error) {
	return t.sidChgTime, nil
}
