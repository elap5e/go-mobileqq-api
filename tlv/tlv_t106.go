package tlv

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"
	"net"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
)

type T106 struct {
	tlv              *TLV
	appID            uint64
	subAppID         uint64
	appClientVersion uint32
	uin              uint64
	serverTime       uint32
	ip               net.IP
	i2               bool
	passwordMD5      [16]byte
	salt             uint64
	username         string
	userA1Key        [16]byte
	isGUIDAvailable  bool
	guid             []byte
	loginType        uint32

	ssoVersion uint32
}

func NewT106(appID, subAppID uint64, appClientVersion uint32, uin uint64, serverTime uint32, ip net.IP, i2 bool, passwordMD5 [16]byte, salt uint64, username string, userA1Key [16]byte, isGUIDAvailable bool, guid []byte, loginType, ssoVersion uint32) *T106 {
	return &T106{
		tlv:              NewTLV(0x0106, 0x0000, nil),
		appID:            appID,
		subAppID:         subAppID,
		appClientVersion: appClientVersion,
		uin:              uin,
		serverTime:       serverTime,
		ip:               ip,
		i2:               i2,
		passwordMD5:      passwordMD5,
		salt:             salt,
		username:         username,
		userA1Key:        userA1Key,
		isGUIDAvailable:  isGUIDAvailable,
		guid:             guid,
		loginType:        loginType,

		ssoVersion: ssoVersion,
	}
}

func (t *T106) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0004)
	v.EncodeUint32(rand.Uint32())
	v.EncodeUint32(t.ssoVersion)
	v.EncodeUint32(uint32(t.appID))
	v.EncodeUint32(t.appClientVersion)
	if t.uin != 0 {
		v.EncodeUint64(t.uin)
	} else {
		v.EncodeUint64(t.salt)
	}
	v.EncodeUint32(t.serverTime)
	v.EncodeRawBytes(t.ip.To4())
	v.EncodeBool(t.i2)
	v.EncodeRawBytes(t.passwordMD5[:])
	v.EncodeRawBytes(t.userA1Key[:])
	v.EncodeUint32(0x00000000)
	v.EncodeBool(t.isGUIDAvailable)
	if len(t.guid) == 0 {
		v.EncodeUint64(rand.Uint64())
		v.EncodeUint64(rand.Uint64())
	} else {
		v.EncodeRawBytes(t.guid)
	}
	v.EncodeUint32(uint32(t.subAppID))
	v.EncodeUint32(t.loginType)
	v.EncodeString(t.username)

	key := append(t.passwordMD5[:], make([]byte, 8)...)
	if t.uin != 0 {
		binary.BigEndian.PutUint64(key[16:], t.uin)
	} else {
		binary.BigEndian.PutUint64(key[16:], t.salt)
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
