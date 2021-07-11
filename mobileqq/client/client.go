package client

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client/config"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/mobileqq/rpc"
)

type Client struct {
	cfg *config.Config
	rpc rpc.Engine

	userSignatures    map[string]*rpc.UserSignature
	userSignaturesMux sync.RWMutex

	channels map[uint64]string
	contacts map[uint64]string

	// message
	messageSeq map[string]*uint32
	syncCookie []byte
}

func NewClient(cfg *config.Config, rpc rpc.Engine) *Client {
	c := &Client{
		cfg: cfg,
		rpc: rpc,
	}
	c.init()
	return c
}

func (c *Client) init() {
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
	c.messageSeq = make(map[string]*uint32)
}

func (c *Client) setMessageSeq(id string, seq uint32) bool {
	if _, ok := c.messageSeq[id]; !ok {
		c.messageSeq[id] = &[]uint32{seq}[0]
	}
	if *c.messageSeq[id] < seq {
		atomic.StoreUint32(c.messageSeq[id], seq)
		return true
	}
	return false
}

func (c *Client) getNextMessageSeq(id string) uint32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if _, ok := c.messageSeq[id]; !ok {
		c.messageSeq[id] = &[]uint32{uint32(r.Int31n(1000)) + 600}[0]
	}
	seq := atomic.AddUint32(c.messageSeq[id], 1)
	if seq > 60000 {
		c.messageSeq[id] = &[]uint32{uint32(r.Int31n(1000)) + 600}[0]
	}
	return seq
}

func (c *Client) Call(
	serviceMethod string,
	c2s *codec.ClientToServerMessage,
	s2c *codec.ServerToClientMessage,
	timeout ...time.Duration,
) error {
	return c.rpc.Call(serviceMethod, c2s, s2c, timeout...)
}

func (c *Client) GetEngine() rpc.Engine {
	return c.rpc
}

func (c *Client) GetNextSeq() uint32 {
	return c.rpc.GetNextSeq()
}

func (c *Client) GetNextMessageSeq(id string) uint32 {
	return c.getNextMessageSeq(id)
}

var clientCtxKey struct{}

func (c *Client) WithClient(ctx context.Context) context.Context {
	return context.WithValue(ctx, clientCtxKey, c)
}

func ForClient(ctx context.Context) *Client {
	return ctx.Value(clientCtxKey).(*Client)
}
