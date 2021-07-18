package client

import (
	"context"
	"database/sql"
	"math/rand"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client/config"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/mobileqq/highway"
	"github.com/elap5e/go-mobileqq-api/mobileqq/rpc"
)

type Client struct {
	mux sync.Mutex

	cfg *config.Config
	db  *sql.DB
	rpc rpc.Engine

	userSignatures    map[string]*rpc.UserSignature
	userSignaturesMux sync.RWMutex

	channels map[int64]*db.Channel
	cmembers map[int64]map[int64]*db.ChannelMember
	contacts map[int64]*db.Contact

	requestSeq int32

	// message
	messageSeq map[string]*int32
	syncCookie []byte
}

func NewClient(cfg *config.Config, rpc rpc.Engine) *Client {
	c := &Client{
		cfg:      cfg,
		rpc:      rpc,
		channels: make(map[int64]*db.Channel),
		cmembers: make(map[int64]map[int64]*db.ChannelMember),
		contacts: make(map[int64]*db.Contact),
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
	c.messageSeq = make(map[string]*int32)
}

func (c *Client) setMessageSeq(id string, seq int32) bool {
	if _, ok := c.messageSeq[id]; !ok {
		c.messageSeq[id] = &[]int32{seq}[0]
	}
	if *c.messageSeq[id] < seq {
		atomic.StoreInt32(c.messageSeq[id], seq)
		return true
	}
	return false
}

func (c *Client) getNextMessageSeq(id string) int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if _, ok := c.messageSeq[id]; !ok {
		c.messageSeq[id] = &[]int32{r.Int31n(1000) + 600}[0]
	}
	seq := atomic.AddInt32(c.messageSeq[id], 1)
	if seq > 60000 {
		c.messageSeq[id] = &[]int32{r.Int31n(1000) + 600}[0]
	}
	return seq
}

func (c *Client) getNextRequestSeq() int32 {
	seq := atomic.AddInt32(&c.requestSeq, 1)
	if seq > 1000000 {
		c.requestSeq = rand.Int31n(100000) + 60000
	}
	return seq
}

func (c *Client) GetCacheByUsernameDir(username string) string {
	return path.Join(c.cfg.CacheDir, "by-username", username)
}

func (c *Client) GetCacheDownloadsDir() string {
	return path.Join(c.cfg.CacheDir, "downloads")
}

func (c *Client) GetHighway(addr, username string) *highway.Highway {
	return highway.NewHighway(addr, username, c.cfg.Client.AppID)
}

func (c *Client) Call(
	serviceMethod string,
	c2s *codec.ClientToServerMessage,
	s2c *codec.ServerToClientMessage,
) error {
	d := time.Now().Add(30 * time.Second)
	return c.rpc.CallWithDeadline(serviceMethod, c2s, s2c, d)
}

func (c *Client) GetEngine() rpc.Engine {
	return c.rpc
}

func (c *Client) GetNextSeq() uint32 {
	return c.rpc.GetNextSeq()
}

func (c *Client) GetNextMessageSeq(id string) int32 {
	return c.getNextMessageSeq(id)
}

func (c *Client) SetDB(db *sql.DB) {
	c.mux.Lock()
	c.db = db
	c.mux.Unlock()
}

var clientCtxKey struct{}

func (c *Client) WithClient(ctx context.Context) context.Context {
	return context.WithValue(ctx, clientCtxKey, c)
}

func ForClient(ctx context.Context) *Client {
	return ctx.Value(clientCtxKey).(*Client)
}
