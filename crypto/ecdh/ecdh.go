package ecdh

import (
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"math/big"
)

var X509Prefix = []byte{
	0x30, 0x59, 0x30, 0x13, 0x06, 0x07, 0x2a, 0x86,
	0x48, 0xce, 0x3d, 0x02, 0x01, 0x06, 0x08, 0x2a,
	0x86, 0x48, 0xce, 0x3d, 0x03, 0x01, 0x07, 0x03,
	0x42, 0x00,
}

type PublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}

func (pub *PublicKey) Bytes() []byte {
	return append(append([]byte{0x04}, pub.X.Bytes()...), pub.Y.Bytes()...)
}

type PrivateKey struct {
	PublicKey
	D []byte
}

func (priv *PrivateKey) Public() *PublicKey {
	return &priv.PublicKey
}

func (priv *PrivateKey) ShareKey(pub *PublicKey) [16]byte {
	sx, _ := priv.PublicKey.Curve.ScalarMult(pub.X, pub.Y, priv.D)
	return md5.Sum(sx.Bytes()[:16])
}

func GenerateKey() (*PrivateKey, error) {
	c := elliptic.P256()
	k, x, y, err := elliptic.GenerateKey(c, rand.Reader)
	if err != nil {
		return nil, err
	}
	priv := &PrivateKey{}
	priv.Curve = c
	priv.PublicKey.X, priv.PublicKey.Y = x, y
	priv.D = k
	return priv, nil
}
