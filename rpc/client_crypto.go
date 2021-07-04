package rpc

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elap5e/go-mobileqq-api/crypto/ecdh"
)

var (
	defaultServerECDHPublicKey = []byte{
		0x04, 0xeb, 0xca, 0x94, 0xd7, 0x33, 0xe3, 0x99,
		0xb2, 0xdb, 0x96, 0xea, 0xcd, 0xd3, 0xf6, 0x9a,
		0x8b, 0xb0, 0xf7, 0x42, 0x24, 0xe2, 0xb4, 0x4e,
		0x33, 0x57, 0x81, 0x22, 0x11, 0xd2, 0xe6, 0x2e,
		0xfb, 0xc9, 0x1b, 0xb5, 0x53, 0x09, 0x8e, 0x25,
		0xe3, 0x3a, 0x79, 0x9a, 0xdc, 0x7f, 0x76, 0xfe,
		0xb2, 0x08, 0xda, 0x7c, 0x65, 0x22, 0xcd, 0xb0,
		0x71, 0x9a, 0x30, 0x51, 0x80, 0xcc, 0x54, 0xa8,
		0x2e,
	}
	defaultServerRSAPublicKey, _ = base64.StdEncoding.DecodeString("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuJTW4abQJXeVdAODw1CamZH4QJZChyT08ribet1Gp0wpSabIgyKFZAOxeArcCbknKyBrRY3FFI9HgY1AyItH8DOUe6ajDEb6c+vrgjgeCiOiCVyum4lI5Fmp38iHKH14xap6xGaXcBccdOZNzGT82sPDM2Oc6QYSZpfs8EO7TYT7KSB2gaHz99RQ4A/Lel1Vw0krk+DescN6TgRCaXjSGn268jD7lOO23x5JS1mavsUJtOZpXkK9GqCGSTCTbCwZhI33CpwdQ2EHLhiP5RaXZCio6lksu+d8sKTWU1eEiEb3cQ7nuZXLYH7leeYFoPtbFV4RicIWp0/YG+RP7rLPCwIDAQAB")
)

func (c *Client) initRandomKey() {
	c.randomKey = [16]byte{}
	rand.Read(c.randomKey[:])
}

func (c *Client) initServerPublicKey() {
	err := c.setServerPublicKey(defaultServerECDHPublicKey, 0x0001)
	if err != nil {
		log.Fatalf(
			"==> [init] failed to init default server public key, error: %s",
			err.Error(),
		)
	}
	if err := c.updateServerPublicKey(); err != nil {
		log.Printf(
			"==> [init] failed to init updated server public key, error: %s",
			err.Error(),
		)
	}
}

func (c *Client) setServerPublicKey(key []byte, ver uint16) error {
	pub, err := x509.ParsePKIXPublicKey(append(ecdh.X509Prefix, key...))
	if err != nil {
		return err
	}
	c.serverPublicKey = &ecdh.PublicKey{
		Curve: pub.(*ecdsa.PublicKey).Curve,
		X:     pub.(*ecdsa.PublicKey).X,
		Y:     pub.(*ecdsa.PublicKey).Y,
	}
	c.serverPublicKeyVersion = ver
	return nil
}

func (c *Client) updateServerPublicKey() error {
	type ServerPublicKey struct {
		QuerySpan         uint32 `json:"QuerySpan"`
		PublicKeyMetaData struct {
			KeyVersion    uint16 `json:"KeyVer"`
			PublicKey     string `json:"PubKey"`
			PublicKeySign string `json:"PubKeySign"`
		} `json:"PubKeyMeta"`
	}
	resp, err := http.Get(
		"https://keyrotate.qq.com/rotate_key?cipher_suite_ver=305&uin=10000",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := ServerPublicKey{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}
	rsaPub, err := x509.ParsePKIXPublicKey(defaultServerRSAPublicKey)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(fmt.Sprintf(
		"305%d%s",
		data.PublicKeyMetaData.KeyVersion,
		data.PublicKeyMetaData.PublicKey,
	)))
	sig, _ := base64.StdEncoding.DecodeString(
		data.PublicKeyMetaData.PublicKeySign,
	)
	if err := rsa.VerifyPKCS1v15(
		rsaPub.(*rsa.PublicKey),
		crypto.SHA256,
		hashed[:],
		sig,
	); err != nil {
		return err
	}
	key, _ := hex.DecodeString(data.PublicKeyMetaData.PublicKey)
	if err := c.setServerPublicKey(
		key,
		data.PublicKeyMetaData.KeyVersion,
	); err != nil {
		return err
	}
	return nil
}

func (c *Client) initPrivateKey() {
	var err error
	if c.privateKey, err = ecdh.GenerateKey(); err != nil {
		log.Fatalf(
			"==> [init] failed to init private key, error: %s",
			err.Error(),
		)
	}
}

func (c *Client) initRandomPassword() {
	c.randomPassword = [16]byte{}
	for i := range c.randomPassword {
		c.randomPassword[i] = byte(
			0x41 + c.rand.Intn(1)*0x20 + c.rand.Intn(26), // TODO: fix
		)
	}
}
