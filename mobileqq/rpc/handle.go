package rpc

import (
	"context"
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
		Str("method", s2c.ServiceMethod).
		Uint32("seq", s2c.Seq).
		Str("uin", s2c.Username).
		Msg("--> [recv] notify ")
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
						Str("method", s2c.ServiceMethod).
						Uint32("seq", s2c.Seq).
						Str("uin", s2c.Username).
						Msg("<-- [send] handled")
					return
				}
			} else {
				log.Debug().
					Str("method", s2c.ServiceMethod).
					Uint32("seq", s2c.Seq).
					Str("uin", s2c.Username).
					Msg("··· [send] handled")
			}
		}
		log.Error().
			Err(err).
			Str("method", s2c.ServiceMethod).
			Uint32("seq", s2c.Seq).
			Str("uin", s2c.Username).
			Msg("··· [send] handled")
	} else {
		log.Debug().
			Str("method", s2c.ServiceMethod).
			Uint32("seq", s2c.Seq).
			Str("uin", s2c.Username).
			Msg("··· [send] ignored")
	}
}

func (e *engine) Register(serviceMethod string, f HandleFunc) error {
	e.handlers[strings.ToLower(serviceMethod)] = f
	return nil
}
