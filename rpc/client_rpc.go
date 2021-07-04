package rpc

import (
	"context"
	"io"
	"log"
	"net/rpc"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type HandleFunc func(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error)

type ClientCall struct {
	ServiceMethod         string
	ClientToServerMessage *codec.ClientToServerMessage
	ServerToClientMessage *codec.ServerToClientMessage
	Error                 error
	Done                  chan *ClientCall
}

func (c *Client) initHandlers() {
	c.handlers = make(map[string]HandleFunc)
	c.handlers[ServiceMethodPushConfigDomain] = c.handlePushConfigDomain
	c.handlers[ServiceMethodPushConfigRequest] = c.handlePushConfigRequest
	c.handlers[ServiceMethodPushMessageNotify] = c.handlePushMessageNotify
	c.handlers[ServiceMethodPushOnlineGroupMessage] = c.handlePushOnlineGroupMessage
	c.handlers[ServiceMethodPushOnlineSIDExpired] = c.handlePushOnlineSIDExpired
}

func (c *Client) call(s2c *codec.ServerToClientMessage) {
	if handleFunc, ok := c.handlers[s2c.ServiceMethod]; ok {
		ctx := c.WithClient(context.Background())
		c2s, err := handleFunc(ctx, s2c)
		if err != nil {
			return
		}
		if c2s != nil {
			c.preprocessC2S(c2s)
			c.c2sMux.Lock()
			defer c.c2sMux.Unlock()
			if err := c.codec.Write(c2s); err != nil {
				return
			}
		}
	}
}

func (c *Client) send(call *ClientCall) {
	c.c2sMux.Lock()
	defer c.c2sMux.Unlock()

	// Register this call.
	c.mux.Lock()
	if c.shutdown || c.closing {
		c.mux.Unlock()
		call.Error = rpc.ErrShutdown
		call.done()
		return
	}
	seq := call.ClientToServerMessage.Seq
	if seq == 0 {
		seq = c.getNextSeq()
		call.ClientToServerMessage.Seq = seq
	}
	c.pending[seq] = call
	c.mux.Unlock()

	// Encode and send the request.
	c.c2s = call.ClientToServerMessage
	c.c2s.ServiceMethod = call.ServiceMethod
	err := c.codec.Write(c.c2s)
	if err != nil {
		c.mux.Lock()
		call = c.pending[seq]
		delete(c.pending, seq)
		c.mux.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (c *Client) revc() {
	var err error
	var s2c codec.ServerToClientMessage
	for err == nil {
		s2c = codec.ServerToClientMessage{}
		err = c.codec.ReadHead(&s2c)
		if err != nil {
			break
		}
		c.preprocessS2C(&s2c)
		err = c.codec.ReadBody(&s2c)
		if err != nil {
			break
		}
		seq := s2c.Seq
		c.mux.Lock()
		call := c.pending[seq]
		delete(c.pending, seq)
		c.mux.Unlock()

		if call != nil {
			ts2c := call.ServerToClientMessage
			codec.CopyServerToClientMessage(ts2c, &s2c)
			call.done()
		} else {
			ts2c := codec.ServerToClientMessage{}
			codec.CopyServerToClientMessage(&ts2c, &s2c)
			go c.call(&ts2c)
		}
	}
	// Terminate pending calls.
	c.c2sMux.Lock()
	c.mux.Lock()
	c.shutdown = true
	closing := c.closing
	if err == io.EOF {
		if closing {
			err = rpc.ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
	c.mux.Unlock()
	c.c2sMux.Unlock()
	if err != io.EOF && !closing {
		log.Println("rpc: client protocol error:", err)
	}
}

func (call *ClientCall) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// We don't want to block here. It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
		log.Println(
			"rpc: discarding Call reply due to insufficient Done chan capacity",
		)
	}
}

func (c *Client) Close() error {
	c.mux.Lock()
	if c.closing {
		c.mux.Unlock()
		return rpc.ErrShutdown
	}
	c.closing = true
	c.mux.Unlock()
	return c.codec.Close()
}

func (c *Client) preprocessC2S(c2s *codec.ClientToServerMessage) {
	sig := c.GetUserSignature(c2s.Username)
	if d2, ok := sig.Tickets["D2"]; ok {
		c2s.UserD2 = d2.Sig
		copy(c2s.UserD2Key[:], d2.Key)
	}
	c2s.FixID = c.cfg.Client.AppID
	c2s.AppID = c.cfg.Client.AppID
	c2s.NetworkType = 0x01 // 0x00: Others; 0x01: Wi-Fi
	c2s.NetIPFamily = 0x03 // 0x00: Others; 0x01: IPv4; 0x02: IPv6; 0x03: Dual
	if a2, ok := sig.Tickets["A2"]; ok {
		c2s.UserA2 = a2.Sig
	}
	c2s.Cookie = sig.Session.Cookie
	c2s.IMEI = c.cfg.Device.IMEI
	c2s.KSID = sig.Session.KSID
	c2s.IMSI = c.cfg.Device.IMSI
	c2s.Revision = c.cfg.Client.Revision
}

func (c *Client) preprocessS2C(s2c *codec.ServerToClientMessage) {
	sig := c.GetUserSignature(s2c.Username)
	if d2, ok := sig.Tickets["D2"]; ok {
		copy(s2c.UserD2Key[:], d2.Key)
	}
}

func (c *Client) Go(
	serviceMethod string,
	c2s *codec.ClientToServerMessage,
	s2c *codec.ServerToClientMessage,
	done chan *ClientCall,
) *ClientCall {
	call := ClientCall{}
	call.ServiceMethod = serviceMethod
	c.preprocessC2S(c2s)
	call.ClientToServerMessage, call.ServerToClientMessage = c2s, s2c
	if done == nil {
		done = make(chan *ClientCall, 10) // buffered.
	} else {
		// If caller passes done != nil, it must arrange that
		// done has enough buffer for the number of simultaneous
		// RPCs that will be using that channel. If the channel
		// is totally unbuffered, it's best not to run at all.
		if cap(done) == 0 {
			log.Panic("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	c.send(&call)
	return &call
}

func (c *Client) Call(
	serviceMethod string,
	c2s *codec.ClientToServerMessage,
	s2c *codec.ServerToClientMessage,
) error {
	call := <-c.Go(serviceMethod, c2s, s2c, make(chan *ClientCall, 1)).Done
	return call.Error
}
