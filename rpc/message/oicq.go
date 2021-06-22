package message

import (
	"context"
	"encoding/binary"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type OICQMessage struct {
	Version       uint16
	ServiceMethod uint16
	Uin           uint64
	EncryptMethod uint8
	RandomKey     [16]byte
	PublicKey     []byte
	ShareKey      [16]byte
	Type          uint16
	TLVs          map[uint16]tlv.TLVCodec
}

func marshalOICQMessageHead(ctx context.Context, msg *OICQMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	buf.EncodeUint8(0x02)
	buf.EncodeUint16(0x0000)
	buf.EncodeUint16(msg.Version)
	buf.EncodeUint16(msg.ServiceMethod)
	buf.EncodeUint16(0x0001)
	buf.EncodeUint32(uint32(msg.Uin))
	buf.EncodeUint8(0x03)
	buf.EncodeUint8(msg.EncryptMethod)
	buf.EncodeUint8(0x00)
	buf.EncodeUint32(0x00000002)
	buf.EncodeUint32(0x00000000)
	buf.EncodeUint32(0x00000000)
	return buf.Bytes(), nil
}

func marshalOICQMessageData(ctx context.Context, msg *OICQMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	buf.EncodeUint16(msg.Type)
	for i := range msg.TLVs {
		if msg.TLVs[i] == nil {
			delete(msg.TLVs, i)
		}
	}
	buf.EncodeUint16(uint16(len(msg.TLVs)))
	for i := range msg.TLVs {
		msg.TLVs[i].Encode(buf)
	}
	return buf.Bytes(), nil
}

func MarshalOICQMessage(ctx context.Context, msg *OICQMessage) ([]byte, error) {
	head, err := marshalOICQMessageHead(ctx, msg)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer([]byte{})
	switch msg.EncryptMethod {
	case 0x07, 0x87: // ECDH
		buf.EncodeUint8(0x02)
		buf.EncodeUint8(0x01)
		buf.EncodeRawBytes(msg.RandomKey[:])
		buf.EncodeUint16(0x0131)
		buf.EncodeBytes(msg.PublicKey)
	case 0x45: // ST
		buf.EncodeUint8(0x01)
		buf.EncodeUint8(0x03)
		buf.EncodeRawBytes(msg.RandomKey[:])
		buf.EncodeUint16(0x0102)
		buf.EncodeUint16(0x0000)
		msg.ShareKey = msg.RandomKey
	}
	data, err := marshalOICQMessageData(ctx, msg)
	if err != nil {
		return nil, err
	}
	buf.EncodeRawBytes(crypto.NewCipher(msg.ShareKey).Encrypt(data))
	buf.EncodeUint8(0x03)
	ret := append(head, buf.Bytes()...)
	binary.BigEndian.PutUint16(ret[1:], uint16(len(ret)))
	return ret, nil
}

func UnmarshalOICQMessage(ctx context.Context, data []byte, msg *OICQMessage) error {
	return nil
}
