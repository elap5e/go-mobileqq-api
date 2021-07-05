package mobileqq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/elap5e/go-mobileqq-api/rpc"
)

type Option struct {
	Config *Config
}

type Client struct {
	cfg Config

	ctx    context.Context    // immutable
	cancel context.CancelFunc // immutable

	addrs    []*net.TCPAddr
	_conn    io.ReadWriteCloser
	_connMux sync.Mutex
	conns    []io.ReadWriteCloser
	connsMux sync.Mutex

	rpc *rpc.Client
}

func NewClient(opts ...Option) *Client {
	cfg := Config{
		RPC: &rpc.Config{
			Client: rpc.NewClientConfig(),
			Device: rpc.NewDeviceConfig(),
		},
	}
	for _, opt := range opts {
		cfg = *opt.Config
	}
	data, _ := json.MarshalIndent(cfg, "", "  ")
	log.Printf("~v~ [init] dump MobileQQ client config:\n%s", string(data))
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}
	client.init()
	return client
}

func (c *Client) init() {
	log.Printf("==> [init] create connection")
	c._connMux.Lock()
	c._conn, _ = c.createConn(context.Background())
	c._connMux.Unlock()
	log.Printf("==> [init] create rpc client")
	c.rpc = rpc.NewClient(c._conn, rpc.Option{Config: c.cfg.RPC})
	c.ctx = c.rpc.WithClient(c.ctx)
}

func (c *Client) runUntilClosed(ctx context.Context) error {
	return nil
}

func (c *Client) Run(
	ctx context.Context,
	f func(ctx context.Context) error,
) error {
	var err error

	defer c.cancel()
	defer func() {
		c.connsMux.Lock()
		defer c.connsMux.Unlock()
		for _, conn := range c.conns {
			closeErr := conn.Close()
			if !errors.Is(closeErr, context.Canceled) {
				err = fmt.Errorf("%v closeErr:%v;", err, closeErr)
			}
		}
	}()

	return c.runUntilClosed(ctx)
}
