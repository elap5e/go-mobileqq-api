package rpc

import (
	"context"
	"io"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elap5e/go-mobileqq-api/crypto/ecdh"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec/tcp"
)

type Client struct {
	cfg Config

	// codec
	codec codec.ClientCodec

	// rpc
	mux sync.Mutex

	c2s    *codec.ClientToServerMessage
	c2sMux sync.Mutex

	seq      uint32
	pending  map[uint32]*ClientCall
	handlers map[string]HandleFunc
	closing  bool
	shutdown bool

	// message
	syncCookie []byte
	syncSeq    map[uint64]*uint32

	// crypto
	rand *rand.Rand

	randomKey      [16]byte
	randomPassword [16]byte

	privateKey             *ecdh.PrivateKey
	serverPublicKey        *ecdh.PublicKey
	serverPublicKeyVersion uint16

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
	c.initPrivateKey()
	c.initServerPublicKey()

	c.initUserSignatures()

	c.initHandlers()
	c.initSync()
}

func (c *Client) getNextSeq() uint32 {
	seq := atomic.AddUint32(&c.seq, 1)
	if seq > 1000000 {
		c.seq = uint32(rand.Int31n(100000)) + 60000
	}
	return seq - 1
}

func (c *Client) initSync() {
	// c.syncCookie = make(map[uint64][]byte)
	c.syncSeq = make(map[uint64]*uint32)
}

// func (c *Client) getSyncCookie(uin uint64) []byte {
// 	return c.syncCookie[uin]
// }

func (c *Client) setSyncSeq(uin uint64, seq uint32) {
	if _, ok := c.syncSeq[uin]; !ok {
		c.syncSeq[uin] = &[]uint32{0}[0]
	}
	if *c.syncSeq[uin] < seq {
		atomic.StoreUint32(c.syncSeq[uin], seq)
	}
}

func (c *Client) getSyncSeq(uin uint64) uint32 {
	if _, ok := c.syncSeq[uin]; !ok {
		c.syncSeq[uin] = &[]uint32{0}[0]
	}
	return *c.syncSeq[uin]
}

func (c *Client) GetNextSyncSeq(uin uint64) uint32 {
	return c.getNextSyncSeq(uin)
}

func (c *Client) getNextSyncSeq(uin uint64) uint32 {
	if _, ok := c.syncSeq[uin]; !ok {
		c.syncSeq[uin] = &[]uint32{0}[0]
	}
	seq := atomic.AddUint32(c.syncSeq[uin], 1)
	if seq > 1000000 {
		c.syncSeq[uin] = &[]uint32{uint32(rand.Int31n(100000)) + 60000}[0]
	}
	return seq - 1
}

func NewClient(conn io.ReadWriteCloser, opts ...Option) *Client {
	return NewClientWithCodec(tcp.NewClientCodec(conn), opts...)
}

func NewClientWithCodec(codec codec.ClientCodec, opts ...Option) *Client {
	cfg := Config{
		Client: NewClientConfig(),
		Device: NewDeviceConfig(),
	}
	for _, opt := range opts {
		cfg = *opt.Config
	}
	// data, _ := json.MarshalIndent(cfg, "", "  ")
	// log.Printf("~v~ [init] dump RPC client config:\n%s", string(data))
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

var s2cCtxKey struct{}

func WithS2C(ctx context.Context, s2c *codec.ServerToClientMessage) context.Context {
	return context.WithValue(ctx, s2cCtxKey, s2c)
}

func ForS2C(ctx context.Context) *codec.ServerToClientMessage {
	return ctx.Value(s2cCtxKey).(*codec.ServerToClientMessage)
}
