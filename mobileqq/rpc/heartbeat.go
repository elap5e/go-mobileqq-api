package rpc

import (
	"time"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (e *engine) HeartbeatAlive(c2s *codec.ClientToServerMessage) error {
	return e.CallWithDeadline(
		"Heartbeat.Alive", c2s, nil, time.Now().Add(2*time.Second),
	)
}

func (e *engine) HeartbeatPing(c2s *codec.ClientToServerMessage) error {
	return e.CallWithDeadline(
		"Heartbeat.Ping", c2s, nil, time.Now().Add(2*time.Second),
	)
}
