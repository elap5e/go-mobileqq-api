package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T202 struct {
	tlv      *TLV
	md5BSSID [16]byte
	ssid     []byte
}

func NewT202(md5BSSID [16]byte, ssid []byte) *T202 {
	return &T202{
		tlv:      NewTLV(0x0202, 0x0000, nil),
		md5BSSID: md5BSSID,
		ssid:     ssid,
	}
}

func (t *T202) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeBytesN(t.md5BSSID[:], 0x0010)
	v.EncodeBytesN(t.ssid, 0x0020)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T202) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	md5BSSID, err := v.DecodeBytes()
	if err != nil {
		return err
	}
	copy(t.md5BSSID[:], md5BSSID)
	if t.ssid, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
