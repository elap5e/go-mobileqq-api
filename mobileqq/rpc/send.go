package rpc

import (
	"github.com/elap5e/go-mobileqq-api/log"
)

func (e *engine) send(call *Call) {
	e.c2sMux.Lock()
	defer e.c2sMux.Unlock()

	// Register this call.
	e.mux.Lock()
	if e.shutdown || e.closing {
		e.mux.Unlock()
		call.Error = ErrShutdown
		call.done()
		return
	}
	seq := call.ClientToServerMessage.Seq
	if seq == 0 {
		seq = e.GetNextSeq()
		call.ClientToServerMessage.Seq = seq
	}
	e.pending[seq] = call
	e.mux.Unlock()

	// Encode and send the request.
	e.c2s = call.ClientToServerMessage
	e.c2s.ServiceMethod = call.ServiceMethod
	err := e.codec.Write(e.c2s)
	if err != nil {
		e.mux.Lock()
		call = e.pending[seq]
		delete(e.pending, seq)
		e.mux.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	} else {
		log.Debug().
			Uint32("@seq", e.c2s.Seq).
			Str("method", e.c2s.ServiceMethod).
			Str("uin", e.c2s.Username).
			Msg("<-- [send]")
	}
}
