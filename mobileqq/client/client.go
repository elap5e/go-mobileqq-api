package client

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/elap5e/go-mobileqq-api/crypto/ecdh"
	"github.com/elap5e/go-mobileqq-api/mobileqq/rpc"
)

type Client struct {
	cfg *Config
	rpc rpc.Engine

	// crypto
	randomKey      [16]byte
	randomPassword [16]byte

	privateKey             *ecdh.PrivateKey
	serverPublicKey        *ecdh.PublicKey
	serverPublicKeyVersion uint16

	userSignatures    map[string]*rpc.UserSignature
	userSignaturesMux sync.RWMutex

	// message
	channels map[uint64]string
	contacts map[uint64]string

	syncCookie []byte
	syncSeq    map[uint64]*uint32

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
}

func NewClient(cfg *Config, rpc rpc.Engine) *Client {
	c := &Client{
		cfg: cfg,
		rpc: rpc,
	}
	c.init()
	return c
}

func (c *Client) init() {
	c.initCrypto()

	c.initUserSignatures()

	c.initHandlers()
	c.initSync()
}

func (c *Client) initHandlers() {
	c.rpc.Register(ServiceMethodConfigPushDomain, c.handleConfigPushDomain)
	c.rpc.Register(ServiceMethodConfigPushRequest, c.handleConfigPushRequest)
	c.rpc.Register(ServiceMethodMessagePushNotify, c.handleMessagePushNotify)
	c.rpc.Register(ServiceMethodMessagePushReaded, c.handleMessagePushReaded)
	c.rpc.Register(ServiceMethodOnlinePushMessageSyncC2C, c.handleOnlinePushMessage)
	c.rpc.Register(ServiceMethodOnlinePushMessageSyncGroup, c.handleOnlinePushMessage)
	c.rpc.Register(ServiceMethodOnlinePushSIDTicketExpired, c.handleOnlinePushSIDTicketExpired)
}

func (c *Client) initSync() {
	c.syncSeq = make(map[uint64]*uint32)
}

func (c *Client) setSyncSeq(uin uint64, seq uint32) bool {
	if _, ok := c.syncSeq[uin]; !ok {
		c.syncSeq[uin] = &[]uint32{seq}[0]
	}
	if *c.syncSeq[uin] < seq {
		atomic.StoreUint32(c.syncSeq[uin], seq)
		return true
	}
	return false
}

func (c *Client) GetNextSyncSeq(uin uint64) uint32 {
	return c.getNextSyncSeq(uin)
}

func (c *Client) getNextSyncSeq(uin uint64) uint32 {
	if _, ok := c.syncSeq[uin]; !ok {
		c.syncSeq[uin] = &[]uint32{uint32(rand.Int31n(100000)) + 60000}[0]
	}
	seq := atomic.AddUint32(c.syncSeq[uin], 1)
	if seq > 1000000 {
		c.syncSeq[uin] = &[]uint32{uint32(rand.Int31n(100000)) + 60000}[0]
	}
	return seq
}

var clientCtxKey struct{}

func (c *Client) WithClient(ctx context.Context) context.Context {
	return context.WithValue(ctx, clientCtxKey, c)
}

func ForClient(ctx context.Context) *Client {
	return ctx.Value(clientCtxKey).(*Client)
}
