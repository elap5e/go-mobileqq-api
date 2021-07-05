package rpc

import (
	"io"

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
		e.lastRecv.Reset(e.interval)
		call := e.pending[seq]
		delete(e.pending, seq)
		e.mux.Unlock()

		if call != nil {
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
	for _, call := range e.pending {
		call.Error = err
		call.done()
	}
	e.mux.Unlock()
	e.c2sMux.Unlock()
	if err != io.EOF && !closing {
		log.Error().Err(err).Msg("rpc: client protocol error")
	}
}
