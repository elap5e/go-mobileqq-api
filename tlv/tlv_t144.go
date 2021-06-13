package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
)

type T144 struct {
	tlv  *TLV
	key  [16]byte
	tlvs []TLVCodec
}

func NewT144(key [16]byte, tlvs ...TLVCodec) *T144 {
	return &T144{
		tlv:  NewTLV(0x0144, 0x0000, nil),
		key:  key,
		tlvs: tlvs,
	}
}

func (t *T144) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(uint16(len(t.tlvs)))
	for i := range t.tlvs {
		t.tlvs[i].Encode(v)
	}
	t.tlv.SetValue(bytes.NewBuffer(crypto.NewCipher(t.key).Encrypt(v.Bytes())))
	t.tlv.Encode(b)
}

func (t *T144) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	_, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}
