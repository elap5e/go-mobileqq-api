package tlv

import (
	"time"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
)

type T400 struct {
	tlv   *TLV
	key   [16]byte
	j     uint64
	bArr2 []byte
	bArr3 []byte
	j2    uint64
	j3    uint64
	bArr4 []byte
}

func NewT400(key [16]byte, j uint64, bArr2, bArr3 []byte, j2, j3 uint64, bArr4 []byte) *T400 {
	return &T400{
		tlv:   NewTLV(0x0400, 0x0000, nil),
		key:   key,
		j:     j,
		bArr2: bArr2,
		bArr3: bArr3,
		j2:    j2,
		j3:    j3,
		bArr4: bArr4,
	}
}

func (t *T400) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	if len(t.bArr2) == 0 {
		t.bArr2 = make([]byte, 16)
	}
	if len(t.bArr3) == 0 {
		t.bArr3 = make([]byte, 16)
	}
	if len(t.bArr4) == 0 {
		t.bArr4 = make([]byte, 8)
	}
	v.EncodeUint16(0x0001)
	v.EncodeUint64(t.j)
	v.EncodeBytes(t.bArr2)
	v.EncodeBytes(t.bArr3)
	v.EncodeUint32(uint32(t.j2))
	v.EncodeUint32(uint32(t.j3))
	v.EncodeUint32(uint32(time.Now().UnixNano() / 1e6))
	v.EncodeBytes(t.bArr4)
	t.tlv.SetValue(bytes.NewBuffer(crypto.NewCipher(t.key).Encrypt(v.Bytes())))
	t.tlv.Encode(b)
}

func (t *T400) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	_, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}
