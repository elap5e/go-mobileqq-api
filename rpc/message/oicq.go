package message

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"

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
	Code          uint8
	TLVs          map[uint16]tlv.TLVCodec
}

func marshalOICQMessageHead(ctx context.Context, msg *OICQMessage) ([]byte, error) {
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

func unmarshalOICQMessageHead(ctx context.Context, buf *bytes.Buffer, msg *OICQMessage) error {
	var err error
	var tmp uint8
	if _, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if tmp, err = buf.DecodeUint8(); err != nil {
		return err
	} else if tmp != 0x02 {
		return fmt.Errorf("unexpected start, got 0x%x", tmp)
	}
	if _, err = buf.DecodeUint16(); err != nil {
		return err
	}
	if msg.Version, err = buf.DecodeUint16(); err != nil {
		return err
	}
	if msg.ServiceMethod, err = buf.DecodeUint16(); err != nil {
		return err
	}
	if _, err = buf.DecodeUint16(); err != nil {
		return err
	}
	var uin uint32
	if uin, err = buf.DecodeUint32(); err != nil {
		return err
	}
	msg.Uin = uint64(uin)
	if _, err = buf.DecodeUint8(); err != nil {
		return err
	}
	if msg.EncryptMethod, err = buf.DecodeUint8(); err != nil {
		return err
	}
	if _, err = buf.DecodeUint8(); err != nil {
		return err
	}
	return nil
}

func unmarshalOICQMessageData(ctx context.Context, buf *bytes.Buffer, msg *OICQMessage) error {
	var err error
	if msg.Type, err = buf.DecodeUint16(); err != nil {
		return err
	}
	if msg.Code, err = buf.DecodeUint8(); err != nil {
		return err
	}
	var l uint16
	if l, err = buf.DecodeUint16(); err != nil {
		return err
	}
	msg.TLVs = map[uint16]tlv.TLVCodec{}
	for i := 0; i < int(l); i++ {
		tlv := tlv.TLV{}
		tlv.Decode(buf)
		msg.TLVs[tlv.GetType()] = &tlv
	}
	return nil
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
		buf.EncodeUint16(0x0001) // TODO: fix here
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
	// for i := range msg.TLVs {
	// 	log.Printf("<-- [send] dump tlv 0x%04x", i)
	// }
	log.Printf("<-- [send] encryptMethod 0x%02x, dump oicq\n%s", msg.EncryptMethod, hex.Dump(data))
	buf.EncodeRawBytes(crypto.NewCipher(msg.ShareKey).Encrypt(data))
	buf.EncodeUint8(0x03)
	ret := append(head, buf.Bytes()...)
	binary.BigEndian.PutUint32(ret[0:], uint32(len(ret)))
	binary.BigEndian.PutUint16(ret[5:], uint16(len(ret)-4))
	return ret, nil
}

func UnmarshalOICQMessage(ctx context.Context, data []byte, msg *OICQMessage) error {
	buf := bytes.NewBuffer(data)
	if err := unmarshalOICQMessageHead(ctx, buf, msg); err != nil {
		return err
	}
	tmp := buf.Bytes()[:buf.Len()-1]
	switch msg.EncryptMethod {
	case 0x00:
	case 0x03:
		msg.ShareKey = msg.RandomKey
	}
	buf = bytes.NewBuffer(crypto.NewCipher(msg.ShareKey).Decrypt(tmp))
	log.Printf("--> [recv] encryptMethod 0x%02x, dump oicq\n%s", msg.EncryptMethod, hex.Dump(buf.Bytes()))
	if err := unmarshalOICQMessageData(ctx, buf, msg); err != nil {
		return err
	}
	for i := range msg.TLVs {
		v := msg.TLVs[i].(*tlv.TLV)
		buf, _ := v.GetValue()
		switch i {
		case 0x000a: // code 0x00000009, message 服务器繁忙，请你稍后再试。
			fallthrough
		case 0x0146: // code 0x00000009, message 登录失败
			code, _ := buf.DecodeUint32()
			mess, _ := buf.DecodeString()
			log.Printf("--> [recv] dump tlv 0x%04x, code 0x%08x, message %s", i, code, mess)
		case 0x0104: // Session
			log.Printf("--> [recv] dump tlv 0x%04x\n%s", i, hex.Dump(buf.Bytes()))
		case 0x0192: // CAPTCHA URL
			log.Printf("--> [recv] dump tlv 0x%04x, url %s", i, string(buf.Bytes()))
		default:
			log.Printf("--> [recv] dump tlv 0x%04x", i)
		}
	}
	return nil
}
