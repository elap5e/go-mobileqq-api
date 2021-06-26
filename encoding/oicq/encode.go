package oicq

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"log"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
)

func Marshal(ctx context.Context, msg *Message) ([]byte, error) {
	head, err := marshalHead(ctx, msg)
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
		buf.EncodeUint16(msg.KeyVersion)
		buf.EncodeBytes(msg.PublicKey)
	case 0x45: // ST
		buf.EncodeUint8(0x01)
		buf.EncodeUint8(0x03)
		buf.EncodeRawBytes(msg.RandomKey[:])
		buf.EncodeUint16(0x0102)
		buf.EncodeUint16(0x0000)
		msg.ShareKey = msg.RandomKey
	}
	data, err := marshalData(ctx, msg)
	if err != nil {
		return nil, err
	}
	// for i := range msg.TLVs {
	// 	log.Printf("<-- [send] dump tlv 0x%04x", i)
	// }
	log.Printf("<-- [send] encryptMethod 0x%02x, dump oicq:\n%s", msg.EncryptMethod, hex.Dump(data))
	buf.EncodeRawBytes(crypto.NewCipher(msg.ShareKey).Encrypt(data))
	buf.EncodeUint8(0x03)
	ret := append(head, buf.Bytes()...)
	binary.BigEndian.PutUint32(ret[0:], uint32(len(ret)))
	binary.BigEndian.PutUint16(ret[5:], uint16(len(ret)-4))
	return ret, nil
}

func marshalHead(ctx context.Context, msg *Message) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	buf.EncodeUint32(0x00000000)
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

func marshalData(ctx context.Context, msg *Message) ([]byte, error) {
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
