package rpc

import (
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (e *engine) HeartbeatAlive(s2c *codec.ServerToClientMessage) error {
	return e.Call("Heartbeat.Alive", &codec.ClientToServerMessage{}, s2c)
}

func (e *engine) HeartbeatPing(s2c *codec.ServerToClientMessage) error {
	return e.Call("Heartbeat.Ping", &codec.ClientToServerMessage{}, s2c)
}
