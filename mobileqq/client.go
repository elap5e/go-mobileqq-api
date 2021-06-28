package mobileqq

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/elap5e/go-mobileqq-api/rpc"
)

type Client struct {
	ctx    context.Context    // immutable
	cancel context.CancelFunc // immutable

	addrs    []*net.TCPAddr
	conn     io.ReadWriteCloser
	connMux  sync.Mutex
	conns    []io.ReadWriteCloser
	connsMux sync.Mutex

	rpc *rpc.Client
}

func NewClient() *Client {
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		ctx:    ctx,
		cancel: cancel,
	}
	client.init()
	return client
}

func (c *Client) init() {
	log.Printf("==> [init] create connection")
	c.connMux.Lock()
	c.conn, _ = c.createConn(context.Background())
	c.connMux.Unlock()
	log.Printf("==> [init] create rpc client")
	c.rpc = rpc.NewClient(c.conn)
	// c.rpc.PushServiceRegister(c.ctx, rpc.NewPushServiceRegisterRequest())
}

func (c *Client) restoreConnection(ctx context.Context) error {
	return nil
}

func (c *Client) runUntilClosed(ctx context.Context) error {
	return nil
}

func (c *Client) Run(ctx context.Context, f func(ctx context.Context) error) error {
	var err error

	defer c.cancel()
	defer func() {
		c.connsMux.Lock()
		defer c.connsMux.Unlock()
		for _, conn := range c.conns {
			if closeErr := conn.Close(); !errors.Is(closeErr, context.Canceled) {
				err = fmt.Errorf("%v closeErr:%v;", err, closeErr)
			}
		}
	}()

	if err = c.restoreConnection(ctx); err != nil {
		return err
	}

	return c.runUntilClosed(ctx)
}

func (c *Client) HeartbeatAlive() error {
	return c.rpc.HeartbeatAlive()
}
