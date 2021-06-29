package rpc

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elap5e/go-mobileqq-api/crypto/ecdh"
)

type ClientCodec interface {
	Encode(msg *ClientToServerMessage) error
	Decode(msg *ServerToClientMessage) error
	DecodeBody(msg *ServerToClientMessage) error

	Close() error
}

type Client struct {
	cfg Config

	// codec
	codec ClientCodec

	// rpc
	c2s    *ClientToServerMessage
	c2sMux sync.Mutex

	mutex    sync.Mutex
	seq      uint32
	pending  map[uint32]*ClientCall
	closing  bool
	shutdown bool

	// random
	rand *rand.Rand

	randomKey      [16]byte
	randomPassword [16]byte

	serverPublicKey        ecdh.PublicKey
	serverPublicKeyVersion uint16
	privateKey             ecdh.PrivateKey

	userSignatures    map[string]*UserSignature
	userSignaturesMux sync.RWMutex

	// tlvs
	t119 []byte
	t172 []byte // from t161
	t173 []byte // from t161
	t17f []byte // from t161
	t106 []byte // from t169
	t10c []byte // from t169
	t16a []byte // from t169
	t145 []byte // from t169
	t174 []byte
	t17b []byte
	t402 []byte
	t403 []byte

	hashedGUID     [16]byte // t401
	loginExtraData []byte   // from t537

	extraData map[uint16][]byte

	randomSeed []byte
	tgtQR      []byte
}

func (c *Client) init() {
	c.initRandomKey()
	c.initRandomPassword()

	c.initServerPublicKey()
	c.initPrivateKey()

	c.initUserSignatures()
}

func (c *Client) initRandomKey() {
	c.randomKey = [16]byte{}
	c.rand.Read(c.randomKey[:])
	log.Printf("--> [init] dump client random key\n%s", hex.Dump(c.randomKey[:]))
}

func (c *Client) initServerPublicKey() {
	if err := c.setServerPublicKey(defaultServerECDHPublicKey, 0x0001); err != nil {
		log.Fatalf("==> [init] failed to init default server public key, error: %s", err.Error())
	}
	if err := c.updateServerPublicKey(); err != nil {
		log.Printf("==> [init] failed to init updated server public key, error: %s", err.Error())
	}
}

func (c *Client) setServerPublicKey(key []byte, ver uint16) error {
	pub, err := x509.ParsePKIXPublicKey(append(ecdh.X509Prefix, key...))
	if err != nil {
		return err
	}
	c.serverPublicKey = ecdh.PublicKey{
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
	resp, err := http.Get("https://keyrotate.qq.com/rotate_key?cipher_suite_ver=305&uin=10000")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := new(ServerPublicKey)
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return err
	}
	rsaPub, err := x509.ParsePKIXPublicKey(defaultServerRSAPublicKey)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(fmt.Sprintf("%d%d%s", 305, data.PublicKeyMetaData.KeyVersion, data.PublicKeyMetaData.PublicKey)))
	sig, _ := base64.StdEncoding.DecodeString(data.PublicKeyMetaData.PublicKeySign)
	if err := rsa.VerifyPKCS1v15(rsaPub.(*rsa.PublicKey), crypto.SHA256, hashed[:], sig); err != nil {
		return err
	}
	key, _ := hex.DecodeString(data.PublicKeyMetaData.PublicKey)
	if err := c.setServerPublicKey(key, data.PublicKeyMetaData.KeyVersion); err != nil {
		return err
	}
	return nil
}

func (c *Client) initPrivateKey() {
	priv, err := ecdh.GenerateKey(c.rand)
	if err != nil {
		log.Fatalf("==> [init] failed to init private key, error: %s", err.Error())
	}
	c.privateKey = *priv
}

func (c *Client) initRandomPassword() {
	c.randomPassword = [16]byte{}
	for i := range c.randomPassword {
		c.randomPassword[i] = byte(0x41 + c.rand.Intn(1)*0x20 + c.rand.Intn(26))
	}
}

func (c *Client) getNextSeq() uint32 {
	seq := atomic.AddUint32(&c.seq, 1)
	if seq > 1000000 {
		c.seq = uint32(rand.Int31n(100000)) + 60000
	}
	return seq - 1
}

func NewClient(conn io.ReadWriteCloser, opts ...Option) *Client {
	return NewClientWithCodec(NewClientCodec(conn), opts...)
}

func NewClientWithCodec(codec ClientCodec, opts ...Option) *Client {
	cfg := Config{
		Client: NewClientConfig(),
		Device: NewDeviceConfig(),
	}
	for _, opt := range opts {
		cfg = *opt.Config
	}
	data, _ := json.MarshalIndent(cfg, "", "  ")
	log.Printf("~v~ [init] dump RPC client config:\n%s", string(data))
	c := &Client{
		cfg:       cfg,
		codec:     codec,
		seq:       uint32(rand.Int31n(100000)) + 60000,
		pending:   make(map[uint32]*ClientCall),
		rand:      rand.New(rand.NewSource(time.Now().Unix())),
		extraData: make(map[uint16][]byte),
	}
	c.init()
	go c.revc()
	return c
}

var clientCtxKey struct{}

func (c *Client) WithClient(ctx context.Context) context.Context {
	return context.WithValue(ctx, clientCtxKey, c)
}

func ForClient(ctx context.Context) *Client {
	return ctx.Value(clientCtxKey).(*Client)
}
