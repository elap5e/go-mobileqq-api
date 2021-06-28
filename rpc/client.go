package rpc

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elap5e/go-mobileqq-api/crypto/ecdh"
)

type ClientCodec interface {
	Encode(msg *ClientToServerMessage) error
	Decode(msg *ServerToClientMessage) error

	Close() error
}

type Client struct {
	// codec
	codec ClientCodec

	// rpc
	c2s    *ClientToServerMessage
	c2sMux sync.Mutex

	mutex    sync.Mutex
	seq      uint32
	pending  map[uint32]*ClientCall
	closing  bool
	shutdown bool

	// random
	rand *rand.Rand

	cookie         [4]byte
	userA1Key      [16]byte
	randomKey      [16]byte
	randomPassword [16]byte

	serverPublicKey        ecdh.PublicKey
	serverPublicKeyVersion uint16
	privateKey             ecdh.PrivateKey

	// tlvs
	t119 []byte
	t172 []byte // from t161
	t173 []byte // from t161
	t17f []byte // from t161
	t106 []byte // from t169
	t10c []byte // from t169
	t16a []byte // from t169
	t145 []byte // from t169
	t174 []byte
	t17b []byte
	t402 []byte
	t403 []byte
	t547 []byte

	session        []byte   // t104
	ksid           []byte   // t108
	hashedGUID     [16]byte // t401
	loginExtraData []byte   // from t537

	// logger
	logger log.Logger
}

func (c *Client) init() {
	c.initLogger()

	c.initCookie()
	c.initUserA1Key()
	c.initRandomKey()
	c.initRandomPassword()

	c.initServerPublicKey()
	c.initPrivateKey()
}

func (c *Client) initLogger() {
	c.logger = log.Logger{}
}

func (c *Client) initCookie() {
	rand.Read(c.cookie[:])
	log.Printf("--> [init] dump cookie\n%s", hex.Dump(c.cookie[:]))
}

func (c *Client) initUserA1Key() {
	rand.Read(c.userA1Key[:])
	log.Printf("--> [init] dump tgtgt key\n%s", hex.Dump(c.userA1Key[:]))
}

func (c *Client) initRandomKey() {
	c.randomKey = [16]byte{}
	c.rand.Read(c.randomKey[:])
	log.Printf("--> [init] dump client random key\n%s", hex.Dump(c.randomKey[:]))
}

func (c *Client) initServerPublicKey() {
	if err := c.setServerPublicKey(defaultServerECDHPublicKey, 0x0001); err != nil {
		c.logger.Fatalf("==> [init] failed to init default server public key, error: %s", err.Error())
	}
	if err := c.updateServerPublicKey(); err != nil {
		c.logger.Printf("==> [init] failed to init updated server public key, error: %s", err.Error())
	}
}

func (c *Client) setServerPublicKey(key []byte, ver uint16) error {
	pub, err := x509.ParsePKIXPublicKey(append(ecdh.X509Prefix, key...))
	if err != nil {
		return err
	}
	c.serverPublicKey = ecdh.PublicKey{
		Curve: pub.(*ecdsa.PublicKey).Curve,
		X:     pub.(*ecdsa.PublicKey).X,
		Y:     pub.(*ecdsa.PublicKey).Y,
	}
	c.serverPublicKeyVersion = ver
	return nil
}

func (c *Client) updateServerPublicKey() error {
	resp, err := http.Get("https://keyrotate.qq.com/rotate_key?cipher_suite_ver=305&uin=10000")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := new(ServerPublicKey)
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return err
	}
	rsaPub, err := x509.ParsePKIXPublicKey(defaultServerRSAPublicKey)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(fmt.Sprintf("%d%d%s", 305, data.PublicKeyMetaData.KeyVersion, data.PublicKeyMetaData.PublicKey)))
	sig, _ := base64.StdEncoding.DecodeString(data.PublicKeyMetaData.PublicKeySign)
	if err := rsa.VerifyPKCS1v15(rsaPub.(*rsa.PublicKey), crypto.SHA256, hashed[:], sig); err != nil {
		return err
	}
	key, _ := hex.DecodeString(data.PublicKeyMetaData.PublicKey)
	if err := c.setServerPublicKey(key, data.PublicKeyMetaData.KeyVersion); err != nil {
		return err
	}
	return nil
}

func (c *Client) initPrivateKey() {
	priv, err := ecdh.GenerateKey(c.rand)
	if err != nil {
		c.logger.Fatalf("==> [init] failed to init private key, error: %s", err.Error())
	}
	c.privateKey = *priv
}

func (c *Client) initRandomPassword() {
	c.randomPassword = [16]byte{}
	for i := range c.randomPassword {
		c.randomPassword[i] = byte(0x41 + c.rand.Intn(1)*0x20 + c.rand.Intn(26))
	}
}

func (c *Client) getNextSeq() uint32 {
	seq := atomic.AddUint32(&c.seq, 1)
	if seq > 1000000 {
		c.seq = uint32(rand.Int31n(100000)) + 60000
	}
	return seq - 1
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

func NewClient(conn io.ReadWriteCloser) *Client {
	return NewClientWithCodec(NewClientCodec(conn))
}

func NewClientWithCodec(codec ClientCodec) *Client {
	c := &Client{
		codec:   codec,
		seq:     uint32(rand.Int31n(100000)) + 60000,
		pending: make(map[uint32]*ClientCall),
		rand:    rand.New(rand.NewSource(time.Now().Unix())),
	}
	c.init()
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
	c2s.AppID = clientAppID
	c2s.Cookie = c.cookie[:]
	c2s.ReserveField = c.ksid
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

func (c *Client) Call(serviceMethod string, c2s *ClientToServerMessage, s2c *ServerToClientMessage) error {
	call := <-c.Go(serviceMethod, c2s, s2c, make(chan *ClientCall, 1)).Done
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
