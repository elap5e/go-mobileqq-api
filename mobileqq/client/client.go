package client

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elap5e/go-mobileqq-api/log"
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

	// user signature
	userSignatures    map[string]*rpc.UserSignature
	userSignaturesMux sync.RWMutex

	// auto status
	autoStatusTimers    map[string]*time.Timer
	autoStatusTimersMux sync.Mutex

	channels map[int64]*db.Channel
	cmembers map[int64]map[int64]*db.ChannelMember
	contacts map[int64]*db.Contact

	requestSeq int32

	// message
	messageSeq map[int64]map[string]*int32
	syncCookie map[int64][]byte
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
	c.initAutoStatusTimers()

	c.initHandlers()
	c.initMessageSeq()
}

func (c *Client) initHandlers() {
	c.rpc.Register(ServiceMethodConfigPushDomain, c.handleConfigPushDomain)
	c.rpc.Register(ServiceMethodConfigPushRequest, c.handleConfigPushRequest)
	c.rpc.Register(ServiceMethodMessagePushNotify, c.handleMessagePushNotify)
	c.rpc.Register(ServiceMethodMessagePushReaded, c.handleMessagePushReaded)
	c.rpc.Register(ServiceMethodOnlinePushMessageSyncC2C, c.handleOnlinePushMessage)
	c.rpc.Register(ServiceMethodOnlinePushMessageSyncGroup, c.handleOnlinePushMessage)
	c.rpc.Register(ServiceMethodOnlinePushSIDTicketExpired, c.handleOnlinePushSIDTicketExpired)
	c.rpc.Register(ServiceMethodOnlinePushRequest, c.handleOnlinePushRequest)
	c.rpc.Register(ServiceMethodQualityTestPushList, c.handleQualityTestPushList)
}

func (c *Client) initMessageSeq() {
	c.messageSeq = make(map[int64]map[string]*int32)
	c.syncCookie = make(map[int64][]byte)
}

func (c *Client) setMessageSeq(peerID, userID, fromID int64, maxSeq int32) bool {
	if _, ok := c.messageSeq[fromID]; !ok {
		c.messageSeq[fromID] = make(map[string]*int32)
	}
	chatID := fmt.Sprintf("@%du%d", peerID, userID)
	if _, ok := c.messageSeq[fromID][chatID]; !ok {
		c.messageSeq[fromID][chatID] = &[]int32{maxSeq}[0]
	}
	if *c.messageSeq[fromID][chatID] < maxSeq {
		atomic.StoreInt32(c.messageSeq[fromID][chatID], maxSeq)
		if c.db != nil {
			err := c.dbInsertMessageSequence(uint64(fromID), &db.MessageSequence{
				PeerID: peerID,
				UserID: userID,
				Type:   0,
				MaxSeq: maxSeq,
			})
			if err != nil {
				log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageSequence")
			}
		}
		return true
	}
	return false
}

func (c *Client) getNextMessageSeq(peerID, userID, fromID int64) int32 {
	if _, ok := c.messageSeq[fromID]; !ok {
		c.messageSeq[fromID] = make(map[string]*int32)
	}
	chatID := fmt.Sprintf("@%du%d", peerID, userID)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := db.MessageSequence{
		PeerID: peerID,
		UserID: userID,
		Type:   0,
		MaxSeq: 0,
	}
	if c.db != nil {
		err := c.dbSelectMessageSequence(uint64(fromID), &ms)
		if err != nil {
			log.Error().Err(err).Msg(">>> [db  ] dbSelectMessageSequence")
		} else {
			c.messageSeq[fromID][chatID] = &ms.MaxSeq
		}
	}
	if _, ok := c.messageSeq[fromID][chatID]; !ok {
		c.messageSeq[fromID][chatID] = &[]int32{r.Int31n(1000) + 600}[0]
	}
	maxSeq := atomic.AddInt32(c.messageSeq[fromID][chatID], 1)
	if maxSeq >= 60000 {
		c.messageSeq[fromID][chatID] = &[]int32{r.Int31n(1000) + 600}[0]
	}
	if c.db != nil {
		err := c.dbInsertMessageSequence(uint64(fromID), &db.MessageSequence{
			PeerID: peerID,
			UserID: userID,
			Type:   0,
			MaxSeq: maxSeq,
		})
		if err != nil {
			log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageSequence")
		}
	}
	return maxSeq
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

func (c *Client) GetNextMessageSeq(peerID, userID, fromID int64) int32 {
	return c.getNextMessageSeq(peerID, userID, fromID)
}

func (c *Client) SetDB(db *sql.DB) {
	c.mux.Lock()
	c.db = db
	if err := c.dbCreateAccountTable(); err != nil {
		log.Fatal().Err(err).
			Msg("failed to operate database")
	}
	c.mux.Unlock()
}

var clientCtxKey struct{}

func (c *Client) WithClient(ctx context.Context) context.Context {
	return context.WithValue(ctx, clientCtxKey, c)
}

func ForClient(ctx context.Context) *Client {
	return ctx.Value(clientCtxKey).(*Client)
}
