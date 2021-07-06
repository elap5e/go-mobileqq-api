package mobileqq

import (
	"context"
	"encoding/json"
	"net"
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

	addrs []net.Addr

	ctx    context.Context
	cancel context.CancelFunc

	ready chan struct{}
}

func NewClient(opt *Options) *Client {
	opt.init()
	data, _ := json.Marshal(opt)
	log.Debug().
		Msgf("··· [dump] Go MobileQQ API client option:%s", string(data))
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		opt:    opt,
		ctx:    ctx,
		cancel: cancel,
	}
	c.init()
	return c
}

func (c *Client) init() {
	c.rpc = rpc.NewEngine(c.opt.Engine)
	c.client = client.NewClient(c.opt.Client.Engine, c.rpc)
	c.ctx = c.client.WithClient(c.ctx)
	c.ready = make(chan struct{}, 1)
}

func (c *Client) reconnectUntilClosed(ctx context.Context) {
	go func() {
		for {
			cfg := c.rpc.GetConfig()
			uris := append(
				connSocketMobileWiFiIPv4Default,
				connSocketMobileWiFiIPv6Default...,
			)
			cfg.Address, _ = c.benchmark(uris)
			c.rpc.SetConfig(cfg)
			c.rpc.Ready(c.ready)
			if err := c.rpc.Start(ctx); err != nil {
				log.Error().
					Err(err).
					Msg("x-x [conn] failed to start rpc engine, retry in 5 seconds...")
				time.Sleep(5 * time.Second)
			} else {
				return
			}
		}
	}()
	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			log.Error().
				Err(err).
				Msg("x-x [conn] canceled by context")
		} else {
			log.Info().
				Msg("··· [conn] canceled by context")
		}
		c.cancel()
	case <-c.ctx.Done():
		if err := ctx.Err(); err != nil {
			log.Error().
				Err(err).
				Msg("x-x [conn] canceled by client context")
		} else {
			log.Info().
				Msg("··· [conn] canceled by client context")
		}
	}
}

func (c *Client) Run(
	ctx context.Context,
	run func(ctx context.Context) error,
) (err error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		c.reconnectUntilClosed(ctx)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		<-c.ready
		err = run(ctx)
		c.cancel()
		wg.Done()
	}()
	wg.Wait()
	return err
}

func (c *Client) GetClient() *client.Client {
	return c.client
}
