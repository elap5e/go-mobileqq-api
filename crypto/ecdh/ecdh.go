package ecdh

import (
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

var X509Prefix, _ = hex.DecodeString("3059301306072a8648ce3d020106082a8648ce3d030107034200")

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
