package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"

	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

var (
	deviceGUID          = md5.Sum(append(defaultDeviceOSBuildID, defaultDeviceMACAddress...)) // []byte("%4;7t>;28<fclient.5*6")
	deviceGUIDFlag      = uint32((1 << 24 & 0xFF000000) | (0 << 8 & 0xFF00))
	deviceIsGUIDFileNil = false
	deviceIsGUIDGenSucc = true
	deviceIsGUIDChanged = false
	deviceDPWD          = []byte{}

	clientVerifyMethod = uint8(0x82) // 0x00, 0x82
	clientRandomKey    = [16]byte{}
)

var ecdh = crypto.NewECDH()

var (
	clientPackageName  = []byte("com.tencent.mobileqq")
	clientVersionName  = []byte("8.8.3")
	clientRevision     = "8.8.3.b2791edc"
	clientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}

	clientAppID      = uint32(0x20030cb2)
	clientBuildTime  = uint64(0x00000000609b85ad)
	clientSDKVersion = "6.0.0.2476"
	clientSSOVersion = uint32(0x00000011)

	clientCodecAppIDDebug   = []byte("736350642")
	clientCodecAppIDRelease = []byte("736350642")
)

var (
	clientImageType  = uint8(0x01)
	clientMiscBitmap = uint32(0x08f7ff7c)
)

func init() {
	deviceDPWD = func(n int) []byte {
		b := make([]byte, n)
		for i := range b {
			b[i] = byte(0x41 + rand.Intn(1)*0x20 + rand.Intn(26))
		}
		return b
	}(16)
	log.Printf("--> [init] dump device dpwd\n%s", hex.Dump(deviceDPWD))
	clientRandomKey = func() [16]byte { var v [16]byte; rand.Read(v[:]); return v }()
	log.Printf("--> [init] dump client random key\n%s", hex.Dump(clientRandomKey[:]))
	log.Printf("==> [init] update ecdh share key")
	if err := ecdh.UpdateShareKey(); err != nil {
		log.Fatalf("x_x [init] failed to init ecdh, error: %s", err.Error())
	}
}

func SetClientForAndroidPhone() {
	clientPackageName = []byte("com.tencent.mobileqq")
	clientVersionName = []byte("8.8.3")
	clientRevision = "8.8.3.b2791edc"
	clientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}

	clientAppID = uint32(0x20030cb2)
	clientBuildTime = uint64(0x00000000609b85ad)
	clientSDKVersion = "6.0.0.2476"
	clientSSOVersion = uint32(0x00000011)

	clientCodecAppIDDebug = []byte("736350642")
	clientCodecAppIDRelease = []byte("736350642")
	tlv.SetSSOVersion(clientSSOVersion)
}

func SetClientForAndroidPad() {
	clientPackageName = []byte("com.tencent.minihd.qq")
	clientVersionName = []byte("5.9.2")
	clientRevision = "5.9.2.3baec0"
	clientSignatureMD5 = [16]byte{0xaa, 0x39, 0x78, 0xf4, 0x1f, 0xd9, 0x6f, 0xf9, 0x91, 0x4a, 0x66, 0x9e, 0x18, 0x64, 0x74, 0xc7}

	clientAppID = uint32(0x2002fdd5)
	clientBuildTime = uint64(0x000000005f1e8730)
	clientSDKVersion = "6.0.0.2433"
	clientSSOVersion = uint32(0x0000000c)

	clientCodecAppIDDebug = []byte("73636270;")
	clientCodecAppIDRelease = []byte("736346857")
	tlv.SetSSOVersion(clientSSOVersion)
}

func SetClientForiPhone() {
	panic("not implement")
}

func SetClientForiPad() {
	panic("not implement")
}

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

	t104 []byte
	t119 []byte
	t174 []byte
	t17b []byte
	t401 [16]byte
	t402 []byte
	t403 []byte
	t547 []byte
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
	log.Printf("--> [init] dump tgtgt key\n%s", hex.Dump(c.tgtgtKey[:]))
	rand.Read(c.cookie[:])
	log.Printf("--> [init] dump cookie\n%s", hex.Dump(c.cookie[:]))
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
