package oicq

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

func Unmarshal(ctx context.Context, data []byte, msg *Message) error {
	buf := bytes.NewBuffer(data)
	if err := unmarshalHead(ctx, buf, msg); err != nil {
		return err
	}
	tmp := buf.Bytes()[:buf.Len()-1]
	switch msg.EncryptMethod {
	case 0x00:
	case 0x03:
		msg.ShareKey = msg.RandomKey
	}
	buf = bytes.NewBuffer(crypto.NewCipher(msg.ShareKey).Decrypt(tmp))
	log.Printf("--> [recv] encryptMethod 0x%02x, dump oicq:\n%s", msg.EncryptMethod, hex.Dump(buf.Bytes()))
	if err := unmarshalData(ctx, buf, msg); err != nil {
		return err
	}
	tlv.DumpTLVs(ctx, msg.TLVs)
	return nil
}

func unmarshalHead(ctx context.Context, buf *bytes.Buffer, msg *Message) error {
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

func unmarshalData(ctx context.Context, buf *bytes.Buffer, msg *Message) error {
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
