package highway

import (
	"google.golang.org/protobuf/proto"
)

func (hw *Highway) write(body []byte) error {
	head, err := proto.Marshal(hw.req)
	if err != nil {
		return err
	}
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	buf.WriteUint8(0x28)
	buf.WriteUint32(uint32(len(head)))
	buf.WriteUint32(uint32(len(body)))
	buf.Write(head)
	buf.Write(body)
	buf.WriteUint8(0x29)
	_, err = hw.conn.Write(buf.Bytes())
	return err
}

func (hw *Highway) send(call *Call) {
	hw.reqMux.Lock()
	defer hw.reqMux.Unlock()

	hw.mux.Lock()
	seq := call.ReqHead.GetBaseHead().GetSeq()
	if seq == 0 {
		seq = hw.getNextSeq()
		call.ReqHead.GetBaseHead().Seq = seq
	}
	hw.pending[seq] = call
	hw.mux.Unlock()

	hw.req = call.ReqHead
	err := hw.write(call.ReqBody)
	if err != nil {
		hw.mux.Lock()
		call := hw.pending[seq]
		delete(hw.pending, seq)
		hw.mux.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}
