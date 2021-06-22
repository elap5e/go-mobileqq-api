package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"math/big"
)

var (
	defaultECDHPublicKey, _       = hex.DecodeString("04edb8906046f5bfbe9abbc5a88b37d70a6006bfbabc1f0cd49dfb33505e63efc5d78ee4e0a4595033b93d02096dcd3190279211f7b4f6785079e19004aa0e03bc")
	defaultECDHShareKey, _        = hex.DecodeString("c129edba736f4909ecc4ab8e010f46a3")
	defaultECDHServerPublicKey, _ = hex.DecodeString("04EBCA94D733E399B2DB96EACDD3F69A8BB0F74224E2B44E3357812211D2E62EFBC91BB553098E25E33A799ADC7F76FEB208DA7C6522CDB0719A305180CC54A82E")
	defaultECDHX509Prefix, _      = hex.DecodeString("3059301306072a8648ce3d020106082a8648ce3d030107034200")
)

type ECDH struct {
	svrPubKey     *ecdsa.PublicKey
	svrKeyVersion uint32
	priKey        []byte

	PublicKey []byte
	ShareKey  [16]byte
}

func NewECDH() *ECDH {
	ecdh := NewECDHDefault()
	if err := ecdh.initShareKey(); err != nil {
		panic("ecdh initShareKey error: " + err.Error())
	}
	return ecdh
}

func NewECDHDefault() *ECDH {
	var pub, err = x509.ParsePKIXPublicKey(append(defaultECDHX509Prefix, defaultECDHServerPublicKey...))
	if err != nil {
		panic("x509 ParsePKIXPublicKey error: " + err.Error())
	}
	ecdh := &ECDH{
		svrPubKey:     pub.(*ecdsa.PublicKey),
		svrKeyVersion: 1,
		PublicKey:     defaultECDHPublicKey,
	}
	copy(ecdh.ShareKey[:], defaultECDHShareKey)
	return ecdh
}

func (c *ECDH) initShareKey() error {
	var prime256v1 = elliptic.P256()
	var x, y *big.Int
	var err error
	if c.priKey, x, y, err = elliptic.GenerateKey(prime256v1, rand.Reader); err != nil {
		return err
	}
	c.PublicKey = append(append([]byte{0x04}, x.Bytes()...), y.Bytes()...)
	var sx, _ = elliptic.P256().ScalarMult(c.svrPubKey.X, c.svrPubKey.Y, c.priKey)
	var shareKey = md5.Sum(sx.Bytes()[:16])
	c.ShareKey = shareKey
	return nil
}
