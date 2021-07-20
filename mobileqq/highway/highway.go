package highway

import (
	"math/rand"
	"net"
	"sync"
	"sync/atomic"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/pb"
)

var bufPool = bytes.NewPool(0)

type Highway struct {
	mux sync.Mutex

	addr string
	conn net.Conn

	req    *pb.Highway_RequestHead
	reqMux sync.Mutex

	seq     uint32
	pending map[uint32]*Call

	uin   string
	appID uint32
}

type Call struct {
	ReqHead *pb.Highway_RequestHead
	ReqBody []byte
	Resp    *pb.Highway_ResponseHead
	Error   error
	Done    chan *Call
}

func (call *Call) done() {
	select {
	case call.Done <- call:
	default:
	}
}

func NewHighway(addr, uin string, appID uint32) *Highway {
	hw := &Highway{
		addr:  addr,
		uin:   uin,
		appID: appID,
	}
	hw.init()
	return hw
}

func (hw *Highway) init() {
	hw.seq = uint32(rand.Int31n(100000)) + 60000
	hw.pending = make(map[uint32]*Call)
}

func (hw *Highway) getNextSeq() uint32 {
	seq := atomic.AddUint32(&hw.seq, 1)
	if seq > 1000000 {
		hw.seq = uint32(rand.Int31n(100000)) + 60000
	}
	return seq
}

func (hw *Highway) Go(
	reqHead *pb.Highway_RequestHead,
	reqBody []byte,
	resp *pb.Highway_ResponseHead,
	done chan *Call,
) *Call {
	call := Call{}
	call.ReqHead, call.ReqBody, call.Resp = reqHead, reqBody, resp
	if done == nil {
		done = make(chan *Call, 10)
	} else {
		if cap(done) == 0 {
			log.Panic().Msg("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	hw.send(&call)
	return &call
}

func (hw *Highway) Call(
	reqHead *pb.Highway_RequestHead,
	reqBody []byte,
	resp *pb.Highway_ResponseHead,
) error {
	call := <-hw.Go(reqHead, reqBody, resp, make(chan *Call, 1)).Done
	return call.Error
}
