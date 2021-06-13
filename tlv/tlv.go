package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type TLV struct {
	t uint16
	l uint16
	v *bytes.Buffer
}

type TLVCodec interface {
	Encode(b *bytes.Buffer)
	Decode(b *bytes.Buffer) error
}

func NewTLV(t uint16, l uint16, v *bytes.Buffer) *TLV {
	return &TLV{t: t, l: l, v: v}
}

func (t *TLV) SetValue(v *bytes.Buffer) {
	t.v = v
}

func (t *TLV) GetValue() (*bytes.Buffer, error) {
	return t.v, nil
}

func (t *TLV) Encode(b *bytes.Buffer) {
	v := t.v.Bytes()
	t.l = uint16(len(v))
	b.EncodeUint16(t.t)
	b.EncodeUint16(t.l)
	b.EncodeRawBytes(v)
}

func (t *TLV) Decode(b *bytes.Buffer) error {
	var err error
	if t.t, err = b.DecodeUint16(); err != nil {
		return err
	}
	var v []byte
	if v, err = b.DecodeBytes(); err != nil {
		return err
	}
	t.l = uint16(len(v))
	t.v = bytes.NewBuffer(v)
	return nil
}
