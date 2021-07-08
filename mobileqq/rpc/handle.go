package rpc

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type HandleFunc func(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error)

func (e *engine) handle(s2c *codec.ServerToClientMessage) {
	log.Debug().
		Uint32("@seq", s2c.Seq).
		Str("@status", "notify").
		Str("method", s2c.ServiceMethod).
		Str("uin", s2c.Username).
		Msg("--> [recv]")
	if handleFunc, ok := e.handlers[strings.ToLower(s2c.ServiceMethod)]; ok {
		ctx := context.Background()
		c2s, err := handleFunc(ctx, s2c)
		if err == nil {
			if c2s != nil {
				e.withContextC2S(c2s)
				e.c2sMux.Lock()
				defer e.c2sMux.Unlock()
				if err = e.codec.Write(c2s); err == nil {
					log.Debug().
						Uint32("@seq", c2s.Seq).
						Str("@status", "handle").
						Str("method", c2s.ServiceMethod).
						Str("uin", c2s.Username).
						Msg("<-- [send]")
					return
				}
			} else {
				log.Debug().
					Uint32("@seq", s2c.Seq).
					Str("@status", "handle").
					Str("method", s2c.ServiceMethod).
					Str("uin", s2c.Username).
					Msg("··· [send]")
				return
			}
		}
		log.Error().Err(err).
			Uint32("@seq", s2c.Seq).
			Str("@status", "handle").
			Str("method", s2c.ServiceMethod).
			Str("uin", s2c.Username).
			Msg("··· [send]")
	} else {
		log.Debug().Msg(">>> [dump]\n" + hex.Dump(s2c.Buffer))
		log.Warn().
			Uint32("@seq", s2c.Seq).
			Str("@status", "ignore").
			Str("method", s2c.ServiceMethod).
			Str("uin", s2c.Username).
			Msg("··· [send]")
	}
}

func (e *engine) Register(serviceMethod string, f HandleFunc) error {
	e.handlers[strings.ToLower(serviceMethod)] = f
	return nil
}
