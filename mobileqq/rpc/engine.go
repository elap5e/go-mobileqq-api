package rpc

import (
	"context"
	"errors"
	_rand "math/rand"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec/tcp"
)

var rand = _rand.New(_rand.NewSource(time.Now().UnixNano()))

var ErrShutdown = errors.New("connection is shut down")

type Engine interface {
	Start(ctx context.Context)
	Close() error
	Error() chan error

	Call(
		serviceMethod string,
		c2s *codec.ClientToServerMessage,
		s2c *codec.ServerToClientMessage,
	) error
	Go(
		serviceMethod string,
		c2s *codec.ClientToServerMessage,
		s2c *codec.ServerToClientMessage,
		done chan *Call,
	) *Call
	Register(serviceMethod string, f HandleFunc) error

	GetConfig() *Config
	GetNextSeq() uint32
	GetUserSignature(username string) *UserSignature

	SetConfig(cfg *Config)
}

type engine struct {
	cfg   *Config
	ctx   context.Context
	err   chan error
	codec codec.ClientCodec
	sigs  map[string]*UserSignature

	mux sync.Mutex

	c2s    *codec.ClientToServerMessage
	c2sMux sync.Mutex

	seq      uint32
	pending  map[uint32]*Call
	handlers map[string]HandleFunc
	closing  bool
	shutdown bool

	// heartbeat
	callback func()
	interval time.Duration
	lastRecv *time.Timer
}

type Call struct {
	ServiceMethod         string
	ClientToServerMessage *codec.ClientToServerMessage
	ServerToClientMessage *codec.ServerToClientMessage
	Error                 error
	Done                  chan *Call
}

func (call *Call) done() {
	if call.Error != nil {
		log.Error().
			Err(call.Error).
			Str("method", call.ServiceMethod).
			Uint32("seq", call.ClientToServerMessage.Seq).
			Str("uin", call.ClientToServerMessage.Username).
			Msg("--> [recv]")
	} else {
		log.Debug().
			Str("method", call.ServiceMethod).
			Uint32("seq", call.ClientToServerMessage.Seq).
			Str("uin", call.ClientToServerMessage.Username).
			Msg("--> [recv]")
	}
	select {
	case call.Done <- call:
		// ok
	default:
		// We don't want to block here. It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
		log.Warn().Msg(
			"rpc: discarding Call reply due to insufficient Done chan capacity",
		)
	}
}

func NewEngine(cfg *Config) Engine {
	e := &engine{
		cfg: cfg,
		err: make(chan error, 1),
	}
	e.init()
	return e
}

func (e *engine) init() {
	e.sigs = make(map[string]*UserSignature)
	e.seq = uint32(rand.Int31n(100000)) + 60000
	e.pending = make(map[uint32]*Call)
	e.handlers = make(map[string]HandleFunc)
}

func (e *engine) withContextC2S(c2s *codec.ClientToServerMessage) {
	sig := e.GetUserSignature(c2s.Username)
	if d2, ok := sig.Tickets["D2"]; ok {
		c2s.UserD2 = d2.Sig
		copy(c2s.UserD2Key[:], d2.Key)
	}
	c2s.FixID = e.cfg.FixID
	c2s.AppID = e.cfg.AppID
	c2s.NetworkType = e.cfg.NetworkType
	c2s.NetIPFamily = e.cfg.NetIPFamily
	if a2, ok := sig.Tickets["A2"]; ok {
		c2s.UserA2 = a2.Sig
	}
	c2s.Cookie = sig.Session.Cookie
	c2s.IMEI = e.cfg.IMEI
	c2s.KSID = sig.Session.KSID
	c2s.IMSI = e.cfg.IMSI
	c2s.Revision = e.cfg.Revision
}

func (e *engine) withContextS2C(s2c *codec.ServerToClientMessage) {
	sig := e.GetUserSignature(s2c.Username)
	if d2, ok := sig.Tickets["D2"]; ok {
		copy(s2c.UserD2Key[:], d2.Key)
	}
}

func (e *engine) Start(ctx context.Context) {
	e.ctx = ctx
	switch strings.ToLower(e.cfg.Network) {
	case "tcp":
		conn, err := net.Dial(e.cfg.Network, e.cfg.Address)
		if err != nil {
			e.err <- err
			return
		}
		log.Info().
			Msgf("<-> [conn] connected to server %s", conn.RemoteAddr().String())
		e.codec = tcp.NewClientCodec(conn)
		e.closing = false
		e.shutdown = false
	}
	go e.recv()
	e.interval = 30 * time.Second
	e.lastRecv = time.AfterFunc(0, func() {
		if err := e.HeartbeatAlive(); err != nil {
			log.Error().
				Err(err).
				Msg("<-x [conn] heartbeat alive")
			e.err <- err
		} else {
			log.Info().
				Msg("<-> [conn] heartbeat alive")
		}
	})
}

func (e *engine) Close() error {
	e.mux.Lock()
	e.lastRecv.Stop()
	if e.closing {
		e.mux.Unlock()
		return ErrShutdown
	}
	e.closing = true
	e.mux.Unlock()
	return e.codec.Close()
}

func (e *engine) Error() chan error {
	return e.err
}

func (e *engine) Go(
	serviceMethod string,
	c2s *codec.ClientToServerMessage,
	s2c *codec.ServerToClientMessage,
	done chan *Call,
) *Call {
	call := Call{}
	call.ServiceMethod = serviceMethod
	e.withContextC2S(c2s)
	call.ClientToServerMessage, call.ServerToClientMessage = c2s, s2c
	if done == nil {
		done = make(chan *Call, 10) // buffered.
	} else {
		// If caller passes done != nil, it must arrange that
		// done has enough buffer for the number of simultaneous
		// RPCs that will be using that channel. If the channel
		// is totally unbuffered, it's best not to run at all.
		if cap(done) == 0 {
			log.Panic().Msg("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	e.send(&call)
	return &call
}

func (e *engine) Call(
	serviceMethod string,
	c2s *codec.ClientToServerMessage,
	s2c *codec.ServerToClientMessage,
) error {
	call := <-e.Go(serviceMethod, c2s, s2c, make(chan *Call, 1)).Done
	return call.Error
}

func (e *engine) GetConfig() *Config {
	return e.cfg
}

func (e *engine) GetNextSeq() uint32 {
	seq := atomic.AddUint32(&e.seq, 1)
	if seq > 1000000 {
		e.seq = uint32(rand.Int31n(100000)) + 60000
	}
	return seq - 1
}

func (e *engine) GetUserSignature(username string) *UserSignature {
	sig, ok := e.sigs[username]
	if !ok {
		sig = &UserSignature{
			Session: &Session{},
		}
	}
	return sig
}

func (e *engine) SetConfig(cfg *Config) {
	e.mux.Lock()
	e.cfg = cfg
	e.mux.Unlock()
}
