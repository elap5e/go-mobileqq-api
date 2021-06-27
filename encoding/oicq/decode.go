package oicq

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

func Unmarshal(ctx context.Context, data []byte, msg *Message) error {
	n, err := checkVaild(data)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data[5 : n-1])
	if err := unmarshalHead(ctx, buf, msg); err != nil {
		return err
	}
	switch msg.EncryptMethod {
	case EncryptMethod0x00:
	case EncryptMethod0x03:
		msg.ShareKey = msg.RandomKey
	}
	buf = bytes.NewBuffer(crypto.NewCipher(msg.ShareKey).Decrypt(buf.Bytes()))
	log.Printf("--> [recv] encryptMethod 0x%02x, dump oicq:\n%s", msg.EncryptMethod, hex.Dump(buf.Bytes()))
	if err := unmarshalData(ctx, buf, msg); err != nil {
		return err
	}
	tlv.DumpTLVs(ctx, msg.TLVs)
	return nil
}

func checkVaild(v []byte) (int, error) {
	n := int(v[0])<<24 | int(v[1])<<16 | int(v[2])<<8 | int(v[3])<<0
	if len(v) < n {
		return 4, io.ErrUnexpectedEOF
	}
	if v[4] != 0x02 {
		return 5, fmt.Errorf("unexpected prefix, got 0x%x", v[4])
	}
	if v[n-1] != 0x03 {
		return n, fmt.Errorf("unexpected suffix, got 0x%x", v[n-1])
	}
	return n, nil
}

func unmarshalHead(ctx context.Context, buf *bytes.Buffer, msg *Message) error {
	var err error
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
	encryptMethod, err := buf.DecodeUint8()
	if err != nil {
		return err
	}
	msg.EncryptMethod = GetEncryptMethod(encryptMethod)
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
