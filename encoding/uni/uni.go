package uni

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/encoding/jce"
)

type Message struct {
	Version     uint16            `jce:",1" json:"Version,omitempty"`
	PacketType  uint8             `jce:",2" json:"PacketType,omitempty"`
	MessageType uint32            `jce:",3" json:"MessageType,omitempty"`
	RequestID   uint32            `jce:",4" json:"RequestID,omitempty"`
	ServantName string            `jce:",5" json:"ServantName,omitempty"`
	FuncName    string            `jce:",6" json:"FuncName,omitempty"`
	Buffer      []byte            `jce:",7" json:"Buffer,omitempty"`
	Timeout     uint32            `jce:",8" json:"Timeout,omitempty"`
	Context     map[string]string `jce:",9" json:"Context,omitempty"`
	Status      map[string]string `jce:",10" json:"Status,omitempty"`
}

type MapBuffer map[string]map[string][]byte

type MapBufferV3 map[string][]byte

func Marshal(
	ctx context.Context,
	msg *Message,
	opts map[string]interface{},
) ([]byte, error) {
	var err error
	switch msg.Version {
	case 0x0001, 0x0002:
		mapBuf := make(MapBuffer)
		for key, opt := range opts {
			tbuf, err := jce.Marshal(opt)
			if err != nil {
				return nil, err
			}
			mapBuf[key][key] = tbuf // TODO: fix
		}
		msg.Buffer, err = jce.Marshal(mapBuf)
		if err != nil {
			return nil, err
		}
	case 0x0003:
		mapBuf := make(MapBufferV3)
		for key, opt := range opts {
			tbuf, err := jce.Marshal(opt)
			if err != nil {
				return nil, err
			}
			mapBuf[key] = tbuf
		}
		msg.Buffer, err = jce.Marshal(mapBuf)
		if err != nil {
			return nil, err
		}
	}
	return jce.Marshal(msg, true)
}

func Unmarshal(
	ctx context.Context,
	data []byte,
	msg *Message,
	opts map[string]interface{},
) error {
	if err := jce.Unmarshal(data, msg, true); err != nil {
		return err
	}
	switch msg.Version {
	case 0x0001, 0x0002:
		mapBuf := make(MapBuffer)
		if err := jce.Unmarshal(msg.Buffer, mapBuf); err != nil {
			return err
		}
		for key, subMapBuf := range mapBuf {
			for _, buf := range subMapBuf {
				if opt, ok := opts[key]; ok {
					if err := jce.Unmarshal(buf, opt); err != nil {
						return err
					}
				}
			}
		}
	case 0x0003:
		mapBuf := make(MapBufferV3)
		if err := jce.Unmarshal(msg.Buffer, mapBuf); err != nil {
			return err
		}
		for key, buf := range mapBuf {
			if opt, ok := opts[key]; ok {
				if err := jce.Unmarshal(buf, opt); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
