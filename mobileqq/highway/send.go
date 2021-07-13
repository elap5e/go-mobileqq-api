package highway

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/pb"
)

func (hw *Highway) send(req, body []byte) (err error) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	buf.WriteUint8(0x28)
	buf.WriteUint32(uint32(len(req)))
	buf.WriteUint32(uint32(len(body)))
	buf.Write(req)
	buf.Write(body)
	buf.WriteUint8(0x29)
	_, err = hw.conn.Write(buf.Bytes())
	return
}

func (hw *Highway) mustSend(req, body []byte) error {
	hw.mux.Lock()
	defer hw.mux.Unlock()

	err := hw.send(req, body)
	if err != nil {
		return err
	}
	head, _, err := hw.recv()
	if err != nil {
		return err
	}

	resp := pb.HighwayResponseHead{}
	err = proto.Unmarshal(head, &resp)
	if err != nil {
		return err
	}

	switch code := resp.GetErrorCode(); code {
	case 0:
		return nil
	case 67:
		return errors.New("invalid upload key")
	case 81:
		return errors.New("checksum not match")
	case 199:
		return errors.New("nil upload key")
	default:
		return fmt.Errorf("not implement %d", code)
	}
}
