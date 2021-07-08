package client

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
	"net/http"

	"github.com/elap5e/go-mobileqq-api/crypto/ecdh"
	"github.com/elap5e/go-mobileqq-api/log"
)

var (
	serverECDHPublicKey, _ = hex.DecodeString("04ebca94d733e399b2db96eacdd3f69a8bb0f74224e2b44e3357812211d2e62efbc91bb553098e25e33a799adc7f76feb208da7c6522cdb0719a305180cc54a82e")
	serverRSAPublicKey, _  = base64.StdEncoding.DecodeString("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuJTW4abQJXeVdAODw1CamZH4QJZChyT08ribet1Gp0wpSabIgyKFZAOxeArcCbknKyBrRY3FFI9HgY1AyItH8DOUe6ajDEb6c+vrgjgeCiOiCVyum4lI5Fmp38iHKH14xap6xGaXcBccdOZNzGT82sPDM2Oc6QYSZpfs8EO7TYT7KSB2gaHz99RQ4A/Lel1Vw0krk+DescN6TgRCaXjSGn268jD7lOO23x5JS1mavsUJtOZpXkK9GqCGSTCTbCwZhI33CpwdQ2EHLhiP5RaXZCio6lksu+d8sKTWU1eEiEb3cQ7nuZXLYH7leeYFoPtbFV4RicIWp0/YG+RP7rLPCwIDAQAB")
)

func (c *Client) initCrypto() {
	c.initRandomKey()
	c.initRandomPassword()
	c.initPrivateKey()
	c.initServerPublicKey()
}

func (c *Client) initRandomKey() {
	c.randomKey = [16]byte{}
	rand.Read(c.randomKey[:])
}

func (c *Client) initRandomPassword() {
	c.randomPassword = [16]byte{}
	rand.Read(c.randomPassword[:])
	strs := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := range c.randomPassword {
		c.randomPassword[i] = strs[c.randomPassword[i]%52]
	}
}

func (c *Client) initPrivateKey() {
	var err error
	if c.privateKey, err = ecdh.GenerateKey(); err != nil {
		log.Fatal().Err(err).
			Msg("··· [init] failed to generate client ECDH private key")
	}
}

func (c *Client) initServerPublicKey() {
	log.Info().Msg("··· [init] updating server ECDH public key...")
	if err := c.setServerPublicKey(serverECDHPublicKey, 0x0001); err != nil {
		log.Fatal().Err(err).
			Msg("··· [init] failed to set default server ECDH public key")
	}
	if err := c.updateServerPublicKey(); err != nil {
		log.Error().Err(err).
			Msg("··· [init] failed to update server ECDH public key")
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
	type serverPublicKey struct {
		QuerySpan         uint32 `json:"QuerySpan"`
		PublicKeyMetadata struct {
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
	data := serverPublicKey{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}
	pub, err := x509.ParsePKIXPublicKey(serverRSAPublicKey)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(fmt.Sprintf(
		"305%d%s",
		data.PublicKeyMetadata.KeyVersion,
		data.PublicKeyMetadata.PublicKey,
	)))
	sig, _ := base64.StdEncoding.DecodeString(
		data.PublicKeyMetadata.PublicKeySign,
	)
	if err := rsa.VerifyPKCS1v15(
		pub.(*rsa.PublicKey),
		crypto.SHA256,
		hashed[:],
		sig,
	); err != nil {
		return err
	}
	key, _ := hex.DecodeString(data.PublicKeyMetadata.PublicKey)
	if err := c.setServerPublicKey(
		key,
		data.PublicKeyMetadata.KeyVersion,
	); err != nil {
		return err
	}
	return nil
}
