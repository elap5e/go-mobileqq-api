package rpc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	_rand "math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec/tcp"
)

var rand = _rand.New(_rand.NewSource(time.Now().UnixNano()))

var (
	ErrClosedByRemote  = errors.New("connection is closed by remote")
	ErrClosedByTimeout = errors.New("connection is closed by timeout")
	ErrShutdown        = errors.New("connection is shut down")
)

type Engine interface {
	Start(ctx context.Context) error
	Ready(ch chan struct{})
	Close() error

	Call(
		serviceMethod string,
		c2s *codec.ClientToServerMessage,
		s2c *codec.ServerToClientMessage,
	) error
	CallWithDeadline(
		serviceMethod string,
		c2s *codec.ClientToServerMessage,
		s2c *codec.ServerToClientMessage,
		d time.Time,
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
	SetServers(list []string)
	SetUserSignature(username string, sig *UserSignature)
}

var engineCtxKey struct{}

func WithEngine(ctx context.Context, e Engine) context.Context {
	return context.WithValue(ctx, engineCtxKey, e)
}

func ForEngine(ctx context.Context) Engine {
	return ctx.Value(engineCtxKey).(Engine)
}

type engine struct {
	cfg *Config
	mux sync.Mutex

	addrs []string
	codec codec.ClientCodec

	sigs map[string]*UserSignature

	c2s    *codec.ClientToServerMessage
	c2sMux sync.Mutex

	seq      uint32
	pending  map[uint32]*Call
	handlers map[string]HandleFunc
	closing  bool
	shutdown bool

	// heartbeat
	interval time.Duration
	watchDog *time.Timer

	// signals
	err   chan error
	ready chan struct{}
}

type Call struct {
	ServiceMethod         string
	ClientToServerMessage *codec.ClientToServerMessage
	ServerToClientMessage *codec.ServerToClientMessage
	Error                 error
	Done                  chan *Call
}

func (call *Call) done() {
	if call.ServerToClientMessage != nil {
		s2c := call.ServerToClientMessage
		if len(s2c.Buffer) > 0 {
			log.Trace().
				Uint32("@seq", s2c.Seq).
				Str("method", s2c.ServiceMethod).
				Str("uin", s2c.Username).
				Msg(">>> [dump] s2c.Buffer:\n" + hex.Dump(s2c.Buffer))
		}
	}
	if call.Error != nil {
		log.Error().
			Err(call.Error).
			Uint32("@seq", call.ClientToServerMessage.Seq).
			Str("method", call.ServiceMethod).
			Str("uin", call.ClientToServerMessage.Username).
			Msg("--> [recv]")
	} else {
		log.Debug().
			Uint32("@seq", call.ClientToServerMessage.Seq).
			Str("method", call.ServiceMethod).
			Str("uin", call.ClientToServerMessage.Username).
			Msg("--> [recv]")
	}
	select {
	case call.Done <- call:
		// ok
	default:
		// We don't want to block here. It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
		log.Warn().
			Err(fmt.Errorf("insufficient Done chan capacity")).
			Uint32("@seq", call.ClientToServerMessage.Seq).
			Str("method", call.ServiceMethod).
			Str("uin", call.ClientToServerMessage.Username).
			Msg("--> [recv]")
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
	e.tcpTesting(getServerList("wifi"))
}

func (e *engine) reset() {
	e.closing = false
	e.shutdown = false
}

func (e *engine) withContextC2S(c2s *codec.ClientToServerMessage) {
	if c2s.Username == "" {
		c2s.Username = "0"
	}
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

func (e *engine) Start(ctx context.Context) error {
	conn, err := net.Dial("tcp", e.addrs[0])
	if err != nil {
		return err
	}
	log.Info().Msg("<-> [conn] connected to server " + e.addrs[0])
	e.codec = tcp.NewClientCodec(conn)

	e.reset()
	go e.recv()

	e.interval = 60 * time.Second
	e.watchDog = time.AfterFunc(0, func() {
		c2s := codec.ClientToServerMessage{}
		if err := e.HeartbeatAlive(&c2s); err != nil {
			log.Error().Err(err).
				Uint32("@seq", c2s.Seq).
				Str("uin", c2s.Username).
				Msg("x-x [conn] heartbeat alive")
		} else {
			log.Info().
				Uint32("@seq", c2s.Seq).
				Str("uin", c2s.Username).
				Msg("<-> [conn] heartbeat alive")
		}
	})

	e.ready <- struct{}{}
	select {
	case err := <-e.err:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (e *engine) Ready(ch chan struct{}) {
	e.ready = ch
}

func (e *engine) Close() error {
	e.mux.Lock()
	e.watchDog.Stop()
	if e.closing {
		e.mux.Unlock()
		return ErrShutdown
	}
	e.closing = true
	e.mux.Unlock()
	return e.codec.Close()
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

func (e *engine) CallWithDeadline(
	serviceMethod string,
	c2s *codec.ClientToServerMessage,
	s2c *codec.ServerToClientMessage,
	d time.Time,
) error {
	var err error
	err = e.codec.SetReadDeadline(d)
	if err != nil {
		return err
	}
	call := <-e.Go(serviceMethod, c2s, s2c, make(chan *Call, 1)).Done
	err = e.codec.SetReadDeadline(time.Time{})
	if err != nil {
		return err
	}
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
	return seq
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

func (e *engine) SetUserSignature(username string, sig *UserSignature) {
	e.sigs[username] = sig
}
