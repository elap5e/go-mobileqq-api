package rpc

import (
	"log"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (c *Client) HeartbeatAlive() error {
	c2s := &codec.ClientToServerMessage{
		Seq:      c.getNextSeq(),
		Username: "0",
		Buffer:   nil,
		Simple:   false,
	}
	s2c := codec.ServerToClientMessage{}
	if err := c.Call("Heartbeat.Alive", c2s, &s2c); err != nil {
		return err
	}
	log.Printf("rpc.HeartbeatAlive success")
	return nil
}
