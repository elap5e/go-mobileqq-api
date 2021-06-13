package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T147 struct {
	tlv   *TLV
	appID uint64
	bArr  []byte
	bArr2 []byte
}

func NewT147(appID uint64, bArr, bArr2 []byte) *T147 {
	return &T147{
		tlv:   NewTLV(0x0147, 0x0000, nil),
		appID: appID,
		bArr:  bArr,
		bArr2: bArr2,
	}
}

func (t *T147) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint32(uint32(t.appID))
	v.EncodeBytesN(t.bArr, 0x0020)
	v.EncodeBytesN(t.bArr2, 0x0020)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T147) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	appID, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.appID = uint64(appID)
	if t.bArr, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.bArr2, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
