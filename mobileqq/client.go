package mobileqq

import (
	"context"
)

type Client struct {
}

func NewClient(appID int, appHash string) *Client {
	client := &Client{}
	client.init()
	return client
}

func (c *Client) init() {}

func (c *Client) Run(ctx context.Context, f func(ctx context.Context) error) (err error) {
	panic("not implement")
}
