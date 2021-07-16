package highway

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/pb"
)

func (hw *Highway) read(resp *pb.HighwayResponseHead, body []byte) error {
	var err error
	p := make([]byte, 1)
	if _, err = hw.conn.Read(p); err != nil {
		return err
	}
	p = make([]byte, 4)
	if _, err = hw.conn.Read(p); err != nil {
		return err
	}
	l1 := int(p[0])<<24 | int(p[1])<<16 | int(p[2])<<8 | int(p[3])<<0
	p = make([]byte, 4)
	if _, err = hw.conn.Read(p); err != nil {
		return err
	}
	l2 := int(p[0])<<24 | int(p[1])<<16 | int(p[2])<<8 | int(p[3])<<0
	head := make([]byte, l1)
	i, n := 0, 0
	for i < l1 {
		n, err = hw.conn.Read(head[i:])
		if err != nil {
			return err
		}
		i += n
	}
	body = make([]byte, l2)
	i, n = 0, 0
	for i < l2 {
		n, err = hw.conn.Read(body[i:])
		if err != nil {
			return err
		}
		i += n
	}
	p = make([]byte, 1)
	if _, err = hw.conn.Read(p); err != nil {
		return err
	}
	return proto.Unmarshal(head, resp)
}

func (hw *Highway) recv() error {
	var err error
	var resp pb.HighwayResponseHead
	var body []byte
	for err == nil {
		resp = pb.HighwayResponseHead{}
		err = hw.read(&resp, body)
		if err != nil {
			break
		}

		seq := resp.GetBaseHead().GetSeq()
		hw.mux.Lock()
		call := hw.pending[seq]
		delete(hw.pending, seq)
		hw.mux.Unlock()

		if call != nil {
			if call.Resp != nil {
				call.Resp.BaseHead = resp.GetBaseHead()
				call.Resp.CacheCost = resp.GetCacheCost()
				call.Resp.ErrorCode = resp.GetErrorCode()
				call.Resp.ExtendInfo = resp.GetExtendInfo()
			}
			switch code := resp.GetErrorCode(); code {
			case 0:
			case 67:
				call.Error = errors.New("invalid upload key")
			case 81:
				call.Error = errors.New("checksum not match")
			case 199:
				call.Error = errors.New("nil upload key")
			default:
				call.Error = fmt.Errorf("not implement %d", code)
			}
			call.done()
		}
	}
	return err
}
