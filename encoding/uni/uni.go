package uni

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"

	"github.com/elap5e/go-mobileqq-api/encoding/jce"
)

type Message struct {
	Version     uint16            `jce:",1"`
	PacketType  uint8             `jce:",2"`
	MessageType uint32            `jce:",3"`
	RequestID   uint32            `jce:",4"`
	ServantName string            `jce:",5"`
	FuncName    string            `jce:",6"`
	Buffer      map[string][]byte `jce:",7"`
	Timeout     uint32            `jce:",8"`
	Context     map[string]string `jce:",9"`
	Status      map[string]string `jce:",10"`
}

func Marshal(
	ctx context.Context,
	msg *Message,
	opts map[string]interface{},
) ([]byte, error) {
	for key, opt := range opts {
		tbuf, err := jce.Marshal(opt)
		if err != nil {
			return nil, err
		}
		msg.Buffer = map[string][]byte{key: tbuf}
	}
	buf, err := jce.Marshal(msg, true)
	if err != nil {
		return nil, err
	}
	data := append(make([]byte, 4), buf...)
	binary.BigEndian.PutUint32(data[0:], uint32(len(data)))
	return data, nil
}

func Unmarshal(
	ctx context.Context,
	data []byte,
	msg *Message,
	opts map[string]interface{},
) error {
	if int(data[0])<<24+int(data[1])<<16+int(data[2])<<8+int(data[3]) > len(data) {
		return io.ErrUnexpectedEOF
	}
	if err := jce.Unmarshal(data[4:], msg, true); err != nil {
		return err
	}
	for key, buf := range msg.Buffer {
		log.Printf("--> [recv] dump jce, key %s\n%s", key, hex.Dump(buf))
		if opt, ok := opts[key]; ok {
			err := jce.Unmarshal(buf, opt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
