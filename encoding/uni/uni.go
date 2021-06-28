package uni

import (
	"context"

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
	return jce.Marshal(msg)
}

func Unmarshal(ctx context.Context, data []byte, msg *Message, opts map[string]interface{}) error {
	return jce.Unmarshal(data, msg)
}
