package rpc

import (
	"fmt"
	"io"
	"os"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (e *engine) recv() {
	var err error
	var s2c codec.ServerToClientMessage
	for err == nil {
		s2c = codec.ServerToClientMessage{}
		err = e.codec.ReadHead(&s2c)
		if err != nil {
			break
		}
		e.withContextS2C(&s2c)
		err = e.codec.ReadBody(&s2c)
		if err != nil {
			break
		}
		seq := s2c.Seq
		e.mux.Lock()
		e.watchDog.Reset(e.interval)
		call := e.pending[seq]
		delete(e.pending, seq)
		e.mux.Unlock()

		if call != nil {
			if s2c.Code != 0 {
				call.Error = fmt.Errorf("%s(%d)", s2c.Message, s2c.Code)
			}
			if ts2c := call.ServerToClientMessage; ts2c != nil {
				codec.CopyServerToClientMessage(ts2c, &s2c)
			}
			call.done()
		} else {
			ts2c := codec.ServerToClientMessage{}
			codec.CopyServerToClientMessage(&ts2c, &s2c)
			go e.handle(&ts2c)
		}
	}

	// Terminate pending calls.
	e.c2sMux.Lock()
	e.mux.Lock()
	e.shutdown = true
	closing := e.closing
	if err == io.EOF {
		if closing {
			err = ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	for key, call := range e.pending {
		delete(e.pending, key)
		call.Error = err
		call.done()
	}
	e.mux.Unlock()
	e.c2sMux.Unlock()
	if err == io.ErrUnexpectedEOF {
		e.err <- ErrClosedByRemote
	} else if os.IsTimeout(err) {
		e.err <- ErrClosedByTimeout
	} else if err != io.EOF && !closing {
		e.err <- err
		log.Error().Err(err).
			Msg("--> [recv] client protocol")
	}
}
