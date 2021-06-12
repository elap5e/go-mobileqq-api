package tlv

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
)

type T106 struct {
	tlv             *TLV
	appID           uint64
	j2              uint64
	i               uint32
	uin             uint64
	bArr            []byte
	ipAddr          []byte
	i2              bool
	bArr3           []byte
	salt            uint64
	bArr4           []byte
	bArr5           []byte
	isGUIDAvailable bool
	guid            []byte
	i4              uint32
}

func NewT106(appID, j2 uint64, i uint32, uin uint64, bArr, ipAddr []byte, i2 bool, bArr3 []byte, salt uint64, bArr4, bArr5 []byte, isGUIDAvailable bool, guid []byte, i4 uint32) *T106 {
	return &T106{
		tlv:             NewTLV(0x0106, 0x0000, nil),
		appID:           appID,
		j2:              j2,
		i:               i,
		uin:             uin,
		bArr:            bArr,
		ipAddr:          ipAddr,
		i2:              i2,
		bArr3:           bArr3,
		salt:            salt,
		bArr4:           bArr4,
		bArr5:           bArr5,
		isGUIDAvailable: isGUIDAvailable,
		guid:            guid,
		i4:              i4,
	}
}

func (t *T106) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0004)
	v.EncodeUint32(rand.Uint32())
	v.EncodeUint32(0x00000011)
	v.EncodeUint32(uint32(t.appID))
	v.EncodeUint32(t.i)
	if t.uin == 0 {
		v.EncodeUint64(t.salt)
	} else {
		v.EncodeUint64(t.uin)
	}
	v.EncodeRawBytes(t.bArr)
	v.EncodeRawBytes(t.ipAddr)
	v.EncodeBool(t.i2)
	v.EncodeRawBytes(t.bArr3)
	v.EncodeRawBytes(t.bArr5)
	v.EncodeUint32(0x00000000)
	v.EncodeBool(t.isGUIDAvailable)
	if len(t.guid) == 0 {
		v.EncodeUint64(rand.Uint64())
		v.EncodeUint64(rand.Uint64())
	} else {
		v.EncodeRawBytes(t.guid)
	}
	v.EncodeUint32(uint32(t.j2))
	v.EncodeUint32(t.i4)
	v.EncodeBytes(t.bArr4)

	key := append(t.bArr3, make([]byte, 8)...)
	if t.salt == 0 {
		binary.BigEndian.PutUint64(key[16:], t.salt)
	} else {
		binary.BigEndian.PutUint64(key[16:], t.uin)
	}
	t.tlv.SetValue(bytes.NewBuffer(crypto.NewCipher(md5.Sum(key)).Encrypt(v.Bytes())))
	t.tlv.Encode(b)
}

func (t *T106) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	_, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}
