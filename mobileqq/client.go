package mobileqq

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
	"github.com/elap5e/go-mobileqq-api/mobileqq/rpc"
)

type Client struct {
	opt    *Options
	rpc    rpc.Engine
	client *client.Client

	ctx    context.Context
	cancel context.CancelFunc

	ready   chan struct{}
	restart chan struct{}
}

func NewClient(opt *Options) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	opt.init()
	c := &Client{
		opt:     opt,
		ctx:     ctx,
		cancel:  cancel,
		ready:   make(chan struct{}, 1),
		restart: make(chan struct{}, 1),
	}
	c.init()
	return c
}

func (c *Client) init() {
	log.Info().Msg("··· [init] Go MobileQQ API (" + PackageVersion + ")")
	p, _ := json.Marshal(c.opt)
	log.Debug().Msg("··· [init] loaded client option:" + string(p))
	c.rpc = rpc.NewEngine(c.opt.Engine)
	c.rpc.Ready(c.ready)
	c.client = client.NewClient(c.opt.Client.Engine, c.rpc)
}

func (c *Client) connect(ctx context.Context) {
	err := c.rpc.Start(ctx)
	if err != nil && err != rpc.ErrShutdown {
		log.Error().Err(err).
			Msg("x-x [conn] failed to start rpc engine, retry in 5 seconds...")
		c.restart <- struct{}{}
		time.Sleep(5 * time.Second)
	} else if err == rpc.ErrShutdown {
		log.Error().Err(err).
			Msg("x-x [conn] rpc engine shut down")
	}
}

func (c *Client) reconnectUntilClosed(ctx context.Context) {
	go func() {
		for {
			c.connect(ctx)
		}
	}()
	<-c.ctx.Done()
}

func (c *Client) runUntilError(
	ctx context.Context,
	run func(ctx context.Context, once bool, restart chan struct{}) error,
) {
	once := true
	go func() {
		for {
			<-c.ready
			if err := run(ctx, once, c.restart); err != nil {
				log.Error().Err(err).
					Msg("x-x [conn] runtime error")
				if !os.IsTimeout(err) {
					c.cancel()
					return
				}
			}
			once = false
		}
	}()
	<-c.ctx.Done()
}

func (c *Client) Run(
	ctx context.Context,
	run func(ctx context.Context, once bool, restart chan struct{}) error,
) (err error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		c.reconnectUntilClosed(ctx)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		c.runUntilError(ctx, run)
		wg.Done()
	}()
	wg.Wait()
	return err
}

func (c *Client) GetClient() *client.Client {
	return c.client
}
