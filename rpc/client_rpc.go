package rpc

import (
	"io"
	"log"
	"net/rpc"
)

type ClientCall struct {
	ServiceMethod         string
	ClientToServerMessage *ClientToServerMessage
	ServerToClientMessage *ServerToClientMessage
	Error                 error
	Done                  chan *ClientCall
}

type ClientToServerMessage struct {
	Version       uint32
	EncryptType   uint8
	EncryptA2     []byte
	EncryptD2     []byte
	EncryptD2Key  [16]byte
	Username      string
	Seq           uint32
	AppID         uint32
	ServiceMethod string
	Cookie        []byte
	KSID          []byte
	ReserveField  []byte
	Buffer        []byte
	Simple        bool

	CodecAppID       uint32
	CodecIMEI        string
	CodecIMSI        string
	CodecNetworkType uint8
	CodecNetIPFamily uint8
	CodecRevision    string
}

type ServerToClientMessage struct {
	Version       uint32
	EncryptType   uint8
	EncryptD2Key  [16]byte
	Username      string
	Seq           uint32
	ReturnCode    uint32
	ServiceMethod string
	Cookie        []byte
	Buffer        []byte
}

func (c *Client) send(call *ClientCall) {
	c.c2sMux.Lock()
	defer c.c2sMux.Unlock()

	// Register this call.
	c.mutex.Lock()
	if c.shutdown || c.closing {
		c.mutex.Unlock()
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
	c.mutex.Unlock()

	// Encode and send the request.
	c.c2s = call.ClientToServerMessage
	c.c2s.ServiceMethod = call.ServiceMethod
	err := c.codec.Encode(c.c2s)
	if err != nil {
		c.mutex.Lock()
		call = c.pending[seq]
		delete(c.pending, seq)
		c.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (c *Client) revc() {
	var err error
	var s2c ServerToClientMessage
	for err == nil {
		s2c = ServerToClientMessage{}
		err = c.codec.Decode(&s2c)
		if err != nil {
			break
		}
		sig := c.GetUserSignature(s2c.Username)
		if d2, ok := sig.Tickets["D2"]; ok {
			copy(s2c.EncryptD2Key[:], d2.Key)
		}
		err = c.codec.DecodeBody(&s2c)
		if err != nil {
			break
		}
		seq := s2c.Seq
		c.mutex.Lock()
		call := c.pending[seq]
		delete(c.pending, seq)
		c.mutex.Unlock()

		if call != nil {
			ts2c := call.ServerToClientMessage
			ts2c.Version = s2c.Version
			ts2c.EncryptType = s2c.EncryptType
			ts2c.Username = s2c.Username
			ts2c.Seq = s2c.Seq
			ts2c.ReturnCode = s2c.ReturnCode
			ts2c.ServiceMethod = s2c.ServiceMethod
			ts2c.Cookie = s2c.Cookie
			ts2c.Buffer = s2c.Buffer
			call.done()
		} else {
			// server notify
			log.Printf("==> [recv] seq 0x%08x, uin %s, method %s, server notify", seq, s2c.Username, s2c.ServiceMethod)
		}
	}
	// Terminate pending calls.
	c.c2sMux.Lock()
	c.mutex.Lock()
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
	c.mutex.Unlock()
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
		log.Println("rpc: discarding Call reply due to insufficient Done chan capacity")
	}
}

func (c *Client) Close() error {
	c.mutex.Lock()
	if c.closing {
		c.mutex.Unlock()
		return rpc.ErrShutdown
	}
	c.closing = true
	c.mutex.Unlock()
	return c.codec.Close()
}

func (c *Client) preprocess(c2s *ClientToServerMessage, s2c *ServerToClientMessage) {
	c2s.AppID = c.cfg.Client.AppID
	sig := c.GetUserSignature(c2s.Username)
	c2s.Cookie = sig.Session.Cookie
	c2s.KSID = sig.Session.KSID
	if a2, ok := sig.Tickets["A2"]; ok {
		c2s.EncryptA2 = a2.Sig
	}
	if d2, ok := sig.Tickets["D2"]; ok {
		c2s.EncryptD2 = d2.Sig
		copy(c2s.EncryptD2Key[:], d2.Key)
	}
	c2s.CodecAppID = c.cfg.Client.AppID
	c2s.CodecIMEI = defaultDeviceIMEI
	c2s.CodecIMSI = defaultDeviceIMSI
	c2s.CodecNetworkType = 0x01 // 0x00: Others; 0x01: Wi-Fi
	c2s.CodecNetIPFamily = 0x03 // 0x00: Others; 0x01: IPv4; 0x02: IPv6; 0x03: Dual
	c2s.CodecRevision = c.cfg.Client.Revision
}

func (c *Client) Go(serviceMethod string, c2s *ClientToServerMessage, s2c *ServerToClientMessage, done chan *ClientCall) *ClientCall {
	call := new(ClientCall)
	call.ServiceMethod = serviceMethod
	c.preprocess(c2s, s2c)
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
	c.send(call)
	return call
}

func (c *Client) Call(serviceMethod string, c2s *ClientToServerMessage, s2c *ServerToClientMessage) error {
	call := <-c.Go(serviceMethod, c2s, s2c, make(chan *ClientCall, 1)).Done
	return call.Error
}

// func (c *Client) Register(
// 	serviceMethod string,
// 	dispatcher func(context.Context, *ServerToClientMessage) error,
// ) error {
// 	http.HandleFunc()
// 	return nil
// }

// func HandleFunc(
// 	serviceMethod string,
// 	handler func(context.Context, *ServerToClientMessage) error,
// ) error {
// 	return nil
// }

// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }
