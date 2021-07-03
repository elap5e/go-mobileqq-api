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
