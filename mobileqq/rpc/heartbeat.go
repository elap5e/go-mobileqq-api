package rpc

import (
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (e *engine) HeartbeatAlive() error {
	return e.Call("Heartbeat.Alive", &codec.ClientToServerMessage{
		Seq:      e.GetNextSeq(),
		Username: "0",
		Buffer:   nil,
		Simple:   false,
	}, nil)
}

func (e *engine) HeartbeatPing() error {
	return e.Call("Heartbeat.Ping", &codec.ClientToServerMessage{
		Seq:      e.GetNextSeq(),
		Username: "0",
		Buffer:   nil,
		Simple:   false,
	}, nil)
}
