package uni

import (
	"context"
	"encoding/hex"
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

func Marshal(ctx context.Context, msg *Message, opts map[string]interface{}) ([]byte, error) {
	for key, opt := range opts {
		buf, err := jce.Marshal(opt)
		if err != nil {
			return nil, err
		}
		msg.Buffer = map[string][]byte{key: buf}
	}
	return jce.Marshal(msg, true)
}

func Unmarshal(ctx context.Context, data []byte, msg *Message, opts map[string]interface{}) error {
	if err := jce.Unmarshal(data, msg, true); err != nil {
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
