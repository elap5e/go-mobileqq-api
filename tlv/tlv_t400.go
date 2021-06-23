package tlv

import (
	"time"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
)

type T400 struct {
	tlv      *TLV
	key      [16]byte
	uin      uint64
	bArr2    []byte
	dpwd     [16]byte
	appID    uint64
	subAppID uint64
	randSeed []byte
}

func NewT400(key [16]byte, uin uint64, bArr2 []byte, dpwd [16]byte, appID, subAppID uint64, randSeed []byte) *T400 {
	return &T400{
		tlv:      NewTLV(0x0400, 0x0000, nil),
		key:      key,
		uin:      uin,
		bArr2:    bArr2,
		dpwd:     dpwd,
		appID:    appID,
		subAppID: subAppID,
		randSeed: randSeed,
	}
}

func (t *T400) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	if len(t.bArr2) == 0 {
		t.bArr2 = make([]byte, 16)
	}
	if len(t.randSeed) == 0 {
		t.randSeed = make([]byte, 8)
	}
	v.EncodeUint16(0x0001)
	v.EncodeUint64(t.uin)
	v.EncodeBytes(t.bArr2)
	v.EncodeBytes(t.dpwd[:])
	v.EncodeUint32(uint32(t.appID))
	v.EncodeUint32(uint32(t.subAppID))
	v.EncodeUint32(uint32(time.Now().UnixNano() / 1e6))
	v.EncodeBytes(t.randSeed)
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
