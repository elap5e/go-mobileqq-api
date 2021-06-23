package mobileqq

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
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
	c.connMux.Lock()
	c.conn, _ = c.createConn(context.Background())
	c.connMux.Unlock()
	c.rpc = rpc.NewClient(c.conn)
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

func (c *Client) Auth(username, password string) error {
	uin, _ := strconv.Atoi(username)
	c.rpc.AuthGetSessionTicketWithPassword(c.ctx, rpc.NewAuthGetSessionTicketWithPasswordRequest(uint64(uin), password))
	return nil
}
