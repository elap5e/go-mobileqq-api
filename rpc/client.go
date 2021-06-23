package rpc

import (
	"encoding/hex"
	"io"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"

	"github.com/elap5e/go-mobileqq-api/crypto"
)

var ecdh = crypto.NewECDH()

type ClientCodec interface {
	Encode(msg *ClientToServerMessage) error
	Decode(msg *ServerToClientMessage) error

	Close() error
}

type Client struct {
	codec ClientCodec

	c2sMutex sync.Mutex
	c2s      *ClientToServerMessage

	mutex    sync.Mutex
	seq      uint32
	pending  map[uint32]*ClientCall
	closing  bool
	shutdown bool

	tgtgtKey [16]byte
	cookie   [4]byte
}

func (c *Client) getNextSeq() uint32 {
	seq := atomic.AddUint32(&c.seq, 1)
	if seq > 1000000 {
		c.seq = uint32(rand.Int31n(100000)) + 60000
	}
	return seq - 1
}

func (c *Client) send(call *ClientCall) {
	c.c2sMutex.Lock()
	defer c.c2sMutex.Unlock()

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
		seq := s2c.Seq
		c.mutex.Lock()
		call := c.pending[seq]
		delete(c.pending, seq)
		c.mutex.Unlock()

		if call != nil {
			call.ServerToClientMessage.Version = s2c.Version
			call.ServerToClientMessage.EncryptType = s2c.EncryptType
			call.ServerToClientMessage.Username = s2c.Username
			call.ServerToClientMessage.Seq = s2c.Seq
			call.ServerToClientMessage.ReturnCode = s2c.ReturnCode
			call.ServerToClientMessage.ServiceMethod = s2c.ServiceMethod
			call.ServerToClientMessage.Cookie = s2c.Cookie
			call.ServerToClientMessage.Buffer = s2c.Buffer
			call.done()
		}
	}
	// Terminate pending calls.
	c.c2sMutex.Lock()
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
	c.c2sMutex.Unlock()
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

func NewClient(conn io.ReadWriteCloser) *Client {
	return NewClientWithCodec(NewClientCodec(conn))
}

func NewClientWithCodec(codec ClientCodec) *Client {
	c := &Client{
		codec:   codec,
		seq:     uint32(rand.Int31n(100000)) + 60000,
		pending: make(map[uint32]*ClientCall),
	}
	rand.Read(c.tgtgtKey[:])
	log.Printf("==> [init] dump tgtgt key\n%s", hex.Dump(c.tgtgtKey[:]))
	rand.Read(c.cookie[:])
	log.Printf("==> [init] dump cookie\n%s", hex.Dump(c.cookie[:]))
	go c.revc()
	return c
}

func Dial(network, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn), nil
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

func (c *Client) Go(serviceMethod string, c2s *ClientToServerMessage, s2c *ServerToClientMessage, done chan *ClientCall) *ClientCall {
	call := new(ClientCall)
	call.ServiceMethod = serviceMethod
	call.ClientToServerMessage = c2s
	call.ServerToClientMessage = s2c
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

func (c *Client) Call(cmd string, c2s *ClientToServerMessage, s2c *ServerToClientMessage) error {
	call := <-c.Go(cmd, c2s, s2c, make(chan *ClientCall, 1)).Done
	return call.Error
}

func (c *Client) HeartbeatAlive() error {
	c2s := &ClientToServerMessage{
		Seq:      c.getNextSeq(),
		Username: "0",
		Buffer:   nil,
		Simple:   false,
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call("Heartbeat.Alive", c2s, s2c); err != nil {
		return err
	}
	return nil
}
