package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T150 struct {
	tlv     *TLV
	bitmap  uint32
	network uint8
}

func NewT150(bitmap uint32, network uint8) *T150 {
	return &T150{
		tlv:     NewTLV(0x0150, 0x0000, nil),
		bitmap:  bitmap,
		network: network,
	}
}

func (t *T150) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint32(t.bitmap)
	v.EncodeUint8(t.network)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T150) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.bitmap, err = v.DecodeUint32(); err != nil {
		return err
	}
	if t.network, err = v.DecodeUint8(); err != nil {
		return err
	}
	return nil
}

func (t *T150) GetBitmap() (uint32, error) {
	return t.bitmap, nil
}

func (t *T150) GetNetwork() (uint8, error) {
	return t.network, nil
}
