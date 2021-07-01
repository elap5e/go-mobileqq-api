package uni

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/elap5e/go-mobileqq-api/encoding/jce"
)

type Message struct {
	Version     uint16            `jce:",1"`
	PacketType  uint8             `jce:",2"`
	MessageType uint32            `jce:",3"`
	RequestID   uint32            `jce:",4"`
	ServantName string            `jce:",5"`
	FuncName    string            `jce:",6"`
	Buffer      []byte            `jce:",7"`
	Timeout     uint32            `jce:",8"`
	Context     map[string]string `jce:",9"`
	Status      map[string]string `jce:",10"`
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
