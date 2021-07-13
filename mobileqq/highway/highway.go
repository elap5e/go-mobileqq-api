package highway

import (
	"math/rand"
	"net"
	"sync"
	"sync/atomic"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

var bufPool = bytes.NewPool(0)

type Highway struct {
	mux sync.Mutex

	addr string
	conn net.Conn

	seq uint32

	uin   string
	appID uint32
}

func NewHighway(addr, uin string, appID uint32) *Highway {
	return &Highway{
		addr:  addr,
		seq:   uint32(rand.Int31n(100000)) + 60000,
		uin:   uin,
		appID: appID,
	}
}

func (hw *Highway) getNextSeq() uint32 {
	seq := atomic.AddUint32(&hw.seq, 1)
	if seq > 1000000 {
		hw.seq = uint32(rand.Int31n(100000)) + 60000
	}
	return seq
}
