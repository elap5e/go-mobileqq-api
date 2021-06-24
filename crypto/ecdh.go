package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

var (
	defaultECDHPublicKey, _       = hex.DecodeString("04edb8906046f5bfbe9abbc5a88b37d70a6006bfbabc1f0cd49dfb33505e63efc5d78ee4e0a4595033b93d02096dcd3190279211f7b4f6785079e19004aa0e03bc")
	defaultECDHShareKey, _        = hex.DecodeString("c129edba736f4909ecc4ab8e010f46a3")
	defaultECDHServerPublicKey, _ = hex.DecodeString("04EBCA94D733E399B2DB96EACDD3F69A8BB0F74224E2B44E3357812211D2E62EFBC91BB553098E25E33A799ADC7F76FEB208DA7C6522CDB0719A305180CC54A82E")
	defaultECDHX509Prefix, _      = hex.DecodeString("3059301306072a8648ce3d020106082a8648ce3d030107034200")
	defaultRSAServerPublicKey, _  = base64.StdEncoding.DecodeString("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuJTW4abQJXeVdAODw1CamZH4QJZChyT08ribet1Gp0wpSabIgyKFZAOxeArcCbknKyBrRY3FFI9HgY1AyItH8DOUe6ajDEb6c+vrgjgeCiOiCVyum4lI5Fmp38iHKH14xap6xGaXcBccdOZNzGT82sPDM2Oc6QYSZpfs8EO7TYT7KSB2gaHz99RQ4A/Lel1Vw0krk+DescN6TgRCaXjSGn268jD7lOO23x5JS1mavsUJtOZpXkK9GqCGSTCTbCwZhI33CpwdQ2EHLhiP5RaXZCio6lksu+d8sKTWU1eEiEb3cQ7nuZXLYH7leeYFoPtbFV4RicIWp0/YG+RP7rLPCwIDAQAB")
)

type ECDH struct {
	svrPubKey *ecdsa.PublicKey
	priKey    []byte

	KeyVersion uint16
	PublicKey  []byte
	ShareKey   [16]byte
}

type ECDHPublicKey struct {
	QuerySpan         uint32 `json:"QuerySpan"`
	PublicKeyMetaData struct {
		KeyVersion    uint16 `json:"KeyVer"`
		PublicKey     string `json:"PubKey"`
		PublicKeySign string `json:"PubKeySign"`
	} `json:"PubKeyMeta"`
}

func NewECDH() *ECDH {
	ecdh := NewECDHDefault()
	if err := ecdh.initShareKey(); err != nil {
		panic("ecdh initShareKey error: " + err.Error())
	}
	return ecdh
}

func NewECDHDefault() *ECDH {
	pub, err := x509.ParsePKIXPublicKey(append(defaultECDHX509Prefix, defaultECDHServerPublicKey...))
	if err != nil {
		panic("x509 ParsePKIXPublicKey error: " + err.Error())
	}
	ecdh := &ECDH{
		svrPubKey:  pub.(*ecdsa.PublicKey),
		KeyVersion: 1,
		PublicKey:  defaultECDHPublicKey,
	}
	copy(ecdh.ShareKey[:], defaultECDHShareKey)
	return ecdh
}

func (c *ECDH) initShareKey() error {
	var x, y *big.Int
	var err error
	prime256v1 := elliptic.P256()
	if c.priKey, x, y, err = elliptic.GenerateKey(prime256v1, rand.Reader); err != nil {
		return err
	}
	c.PublicKey = append(append([]byte{0x04}, x.Bytes()...), y.Bytes()...)
	sx, _ := elliptic.P256().ScalarMult(c.svrPubKey.X, c.svrPubKey.Y, c.priKey)
	c.ShareKey = md5.Sum(sx.Bytes()[:16])
	return nil
}

func (c *ECDH) UpdateShareKey() error {
	resp, err := http.Get("https://keyrotate.qq.com/rotate_key?cipher_suite_ver=305&uin=10000")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := ECDHPublicKey{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}
	rsaPub, err := x509.ParsePKIXPublicKey(defaultRSAServerPublicKey)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(fmt.Sprintf("%d%d%s", 305, data.PublicKeyMetaData.KeyVersion, data.PublicKeyMetaData.PublicKey)))
	sig, _ := base64.StdEncoding.DecodeString(data.PublicKeyMetaData.PublicKeySign)
	if err := rsa.VerifyPKCS1v15(rsaPub.(*rsa.PublicKey), crypto.SHA256, hashed[:], sig); err != nil {
		return err
	}
	key, _ := hex.DecodeString(data.PublicKeyMetaData.PublicKey)
	pub, err := x509.ParsePKIXPublicKey(append(defaultECDHX509Prefix, key...))
	if err != nil {
		panic("x509 ParsePKIXPublicKey error: " + err.Error())
	}
	c.svrPubKey = pub.(*ecdsa.PublicKey)
	c.KeyVersion = data.PublicKeyMetaData.KeyVersion
	sx, _ := elliptic.P256().ScalarMult(c.svrPubKey.X, c.svrPubKey.Y, c.priKey)
	c.ShareKey = md5.Sum(sx.Bytes()[:16])
	return nil
}
