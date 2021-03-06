package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T109 struct {
	tlv          *TLV
	osBuildIDMD5 [16]byte
}

func NewT109(osBuildIDMD5 [16]byte) *T109 {
	return &T109{
		tlv:          NewTLV(0x0109, 0x0000, nil),
		osBuildIDMD5: osBuildIDMD5,
	}
}

func (t *T109) Encode(b *bytes.Buffer) {
	t.tlv.SetValue(bytes.NewBuffer(t.osBuildIDMD5[:]))
	t.tlv.Encode(b)
}

func (t *T109) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	copy(t.osBuildIDMD5[:], v.Bytes())
	return nil
}
