package uni

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/encoding/jce"
)

type Message struct {
	Version     uint16                 `jce:",1"`
	PacketType  uint8                  `jce:",2"`
	MessageType uint32                 `jce:",3"`
	RequestID   uint32                 `jce:",4"`
	ServantName string                 `jce:",5"`
	FuncName    string                 `jce:",6"`
	Buffer      map[string]interface{} `jce:",7"`
	Timeout     uint32                 `jce:",8"`
	Context     map[string]string      `jce:",9"`
	Status      map[string]string      `jce:",10"`
}

func Marshal(ctx context.Context, msg *Message) ([]byte, error) {
	return jce.Marshal(msg)
}

func Unmarshal(ctx context.Context, data []byte, msg *Message) error {
	return jce.Unmarshal(data, msg)
}
