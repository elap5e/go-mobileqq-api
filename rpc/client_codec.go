package rpc

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
)

type clientCodec struct {
	conn io.ReadWriteCloser

	buf *bytes.Buffer
}

func NewClientCodec(conn io.ReadWriteCloser) ClientCodec {
	return &clientCodec{conn: conn}
}

func (c *clientCodec) encodeHead(msg *ClientToServerMessage) ([]byte, error) {
	if msg.Version != 0x0000000a && msg.Version != 0x0000000b {
		return nil, fmt.Errorf(
			"failed to encode head, version 0x%x",
			msg.Version,
		)
	}
	buf := bytes.NewBuffer([]byte{})
	buf.EncodeUint32(0x00000000)
	buf.EncodeUint32(msg.Version)
	buf.EncodeUint8(msg.EncryptType)
	switch msg.Version {
	case 0x0000000a:
		if msg.EncryptType == 0x01 {
			buf.EncodeUint32(uint32(len(msg.EncryptD2) + 4))
			buf.EncodeRawBytes(msg.EncryptD2)
		} else {
			buf.EncodeUint32(0x00000004)
		}
	case 0x0000000b:
		buf.EncodeUint32(msg.Seq)
	}
	buf.EncodeUint8(0x00)
	buf.EncodeUint32(uint32(len(msg.Username) + 4))
	buf.EncodeRawString(msg.Username)
	return buf.Bytes(), nil
}

func (c *clientCodec) encodeBody(msg *ClientToServerMessage) ([]byte, error) {
	if msg.Version != 0x0000000a && msg.Version != 0x0000000b {
		return nil, fmt.Errorf(
			"failed to encode data, version 0x%x",
			msg.Version,
		)
	}
	buf := bytes.NewBuffer([]byte{})
	buf.EncodeUint32(0x00000000)
	if msg.Version == 0x0000000a {
		buf.EncodeUint32(msg.Seq)
		buf.EncodeUint32(msg.CodecAppID)
		buf.EncodeUint32(msg.AppID)
		{
			tmp := make([]byte, 12)
			tmp[0x0] = msg.CodecNetworkType
			tmp[0xa] = msg.CodecNetIPFamily
			buf.EncodeRawBytes(tmp)
		}
		buf.EncodeUint32(uint32(len(msg.EncryptA2) + 4))
		buf.EncodeRawBytes(msg.EncryptA2)
	}
	buf.EncodeUint32(uint32(len(msg.ServiceMethod) + 4))
	buf.EncodeRawString(msg.ServiceMethod)
	buf.EncodeUint32(uint32(len(msg.Cookie) + 4))
	buf.EncodeRawBytes(msg.Cookie)
	if msg.Version == 0x0000000a {
		buf.EncodeUint32(uint32(len(msg.CodecIMEI) + 4))
		buf.EncodeRawString(msg.CodecIMEI)
		buf.EncodeUint32(uint32(len(msg.KSID) + 4))
		buf.EncodeRawBytes(msg.KSID)
		{
			tmp := "" + "|" + msg.CodecIMSI + "|A" + msg.CodecRevision
			buf.EncodeUint16(uint16(len(tmp) + 2))
			buf.EncodeRawString(tmp)
		}
	}
	buf.EncodeUint32(uint32(len(msg.ReserveField) + 4))
	buf.EncodeRawBytes(msg.ReserveField)
	ret := buf.Bytes()
	binary.BigEndian.PutUint32(ret[0:], uint32(len(ret)))
	if len(msg.Buffer) != 0 {
		ret = append(ret, msg.Buffer...)
	} else {
		tmp := make([]byte, 4)
		binary.BigEndian.PutUint32(tmp, 0x00000004)
		ret = append(ret, tmp...)
	}
	return ret, nil
}

func (c *clientCodec) decodeHead(
	buf *bytes.Buffer,
	msg *ServerToClientMessage,
) error {
	var err error
	if _, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.Version, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.Version != 0x0000000a && msg.Version != 0x0000000b {
		return fmt.Errorf("failed to decode head, version 0x%x", msg.Version)
	}
	if msg.EncryptType, err = buf.DecodeUint8(); err != nil {
		return err
	}
	if _, err = buf.DecodeUint8(); err != nil {
		return err
	} // 0x00
	l, err := buf.DecodeUint32()
	if err != nil {
		return err
	}
	if msg.Username, err = buf.DecodeStringN(uint16(l - 4)); err != nil {
		return err
	}
	return nil
}

func (c *clientCodec) decodeBody(
	buf *bytes.Buffer,
	msg *ServerToClientMessage,
) error {
	if msg.Version != 0x0000000a && msg.Version != 0x0000000b {
		return fmt.Errorf("failed to encode head, version 0x%x", msg.Version)
	}
	var err error
	var n uint32
	if n, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.Seq, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.Code, err = buf.DecodeUint32(); err != nil {
		return err
	}
	var l uint32
	if l, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.Message, err = buf.DecodeStringN(uint16(l - 4)); err != nil {
		return err
	}
	if l, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.ServiceMethod, err = buf.DecodeStringN(uint16(l - 4)); err != nil {
		return err
	}
	if l, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.Cookie, err = buf.DecodeBytesN(uint16(l - 4)); err != nil {
		return err
	}
	if msg.Flag, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if buf.Index() < int(n) {
		if l, err = buf.DecodeUint32(); err != nil {
			return err
		}
		if msg.ReserveField, err = buf.DecodeBytesN(uint16(l - 4)); err != nil {
			return err
		}
	}
	msg.Buffer = buf.Bytes()
	return nil
}

func (c *clientCodec) Encode(msg *ClientToServerMessage) error {
	var err error
	if !msg.Simple {
		msg.Version = 0x0000000a
	} else {
		msg.Version = 0x0000000b
	}
	// log.Printf(
	// 	"  < [send] seq 0x%08x, uin %s, method %s, dump buff:\n%s",
	// 	msg.Seq, msg.Username, msg.ServiceMethod,
	// 	hex.Dump(msg.Buffer),
	// )
	var data []byte
	data, err = c.encodeBody(msg)
	if err != nil {
		return err
	}
	// log.Printf(
	// 	" <- [send] seq 0x%08x, uin %s, method %s, dump data:\n%s",
	// 	msg.Seq, msg.Username, msg.ServiceMethod,
	// 	hex.Dump(data),
	// )
	method := strings.ToLower(msg.ServiceMethod)
	if method == "heartbeat.ping" ||
		method == "heartbeat.alive" ||
		method == "client.correcttime" {
		msg.EncryptType = 0x00
	} else {
		cipher := crypto.NewCipher([16]byte{})
		if len(msg.EncryptD2) == 0 ||
			method == "login.auth" ||
			method == "login.chguin" ||
			method == "grayuinpro.check" ||
			method == "wtlogin.login" ||
			method == "wtlogin.name2uin" ||
			method == "wtlogin.exchange_emp" ||
			method == "wtlogin.trans_emp" ||
			method == "account.requestverifywtlogin_emp" ||
			method == "account.requestrebindmblwtLogin_emp" ||
			method == "connauthsvr.get_app_info_emp" ||
			method == "connauthsvr.get_auth_api_list_emp" ||
			method == "connauthsvr.sdk_auth_api_emp" ||
			method == "qqconnectlogin.pre_auth_emp" ||
			method == "qqconnectlogin.auth_emp" {
			msg.EncryptType = 0x02
		} else {
			cipher.SetKey(msg.EncryptD2Key)
			msg.EncryptType = 0x01
		}
		data = cipher.Encrypt(data)
	}
	var head []byte
	head, err = c.encodeHead(msg)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(head[0:], uint32(len(head)+len(data)))
	if _, err = c.conn.Write(head); err != nil {
		return err
	}
	if _, err = c.conn.Write(data); err != nil {
		return err
	}
	// log.Printf(
	// 	" <= [send] seq 0x%08x, uin %s, method %s, dump send:\n%s",
	// 	msg.Seq, msg.Username, msg.ServiceMethod,
	// 	hex.Dump(append(head, data...)),
	// )
	log.Printf(
		"<== [send] seq 0x%08x, uin %s, method %s",
		msg.Seq, msg.Username, msg.ServiceMethod,
	)
	return nil
}

func (c *clientCodec) Decode(msg *ServerToClientMessage) error {
	var err error
	v := make([]byte, 4)
	if err = c.loopRead(v); err != nil {
		return err
	}
	l := uint32(v[0])<<24 | uint32(v[1])<<16 | uint32(v[2])<<8 | uint32(v[3])<<0
	v = append(v, make([]byte, l-4)...)
	if err = c.loopRead(v[4:]); err != nil {
		return err
	}
	c.buf = bytes.NewBuffer(v)
	if err = c.decodeHead(c.buf, msg); err != nil {
		log.Printf(
			">   [recv] seq 0xffffffff, uin %s, method Unknown, error %v, dump recv:\n%s",
			msg.Username, err,
			hex.Dump(v),
		)
		return err
	}
	// log.Printf(
	// 	">   [recv] seq 0xffffffff, uin %s, method Unknown, dump recv:\n%s",
	// 	msg.Username,
	// 	hex.Dump(v),
	// )
	return nil
}

func (c *clientCodec) DecodeBody(msg *ServerToClientMessage) error {
	switch msg.EncryptType {
	case 0x00:
	case 0x01:
		c.buf = bytes.NewBuffer(
			crypto.NewCipher(msg.EncryptD2Key).Decrypt(c.buf.Bytes()),
		)
	case 0x02:
		c.buf = bytes.NewBuffer(
			crypto.NewCipher([16]byte{}).Decrypt(c.buf.Bytes()),
		)
	default:
		return fmt.Errorf(
			"failed to decode data, encrypt type 0x%x",
			msg.EncryptType,
		)
	}
	v := c.buf.Bytes()
	if err := c.decodeBody(c.buf, msg); err != nil {
		log.Printf(
			"->  [recv] seq 0x%08x, uin %s, method %s, error %v, dump data:\n%s",
			msg.Seq, msg.Username, msg.ServiceMethod, err,
			hex.Dump(v),
		)
		return err
	}
	// log.Printf(
	// 	"->  [recv] seq 0x%08x, uin %s, method %s, dump data:\n%s",
	// 	msg.Seq, msg.Username, msg.ServiceMethod,
	// 	hex.Dump(v),
	// )
	log.Printf(
		"=>  [recv] seq 0x%08x, uin %s, method %s, dump buff:\n%s",
		msg.Seq, msg.Username, msg.ServiceMethod,
		hex.Dump(msg.Buffer),
	)
	log.Printf(
		"==> [recv] seq 0x%08x, uin %s, method %s, code 0x%08x, message %s",
		msg.Seq, msg.Username, msg.ServiceMethod, msg.Code, msg.Message,
	)
	return nil
}

func (c *clientCodec) Close() error {
	return c.conn.Close()
}

func (c *clientCodec) loopRead(v []byte) error {
	l := len(v)
	i := 0
	for i < l {
		n, err := c.conn.Read(v[i:])
		if err != nil {
			return err
		}
		i += n
	}
	return nil
}
