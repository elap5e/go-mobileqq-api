package rpc

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
)

type clientCodec struct {
	conn io.ReadWriteCloser

	appID       uint32
	networkType uint8
	netIPFamily uint8
	imei        string
	imsi        string
	revision    string

	authData map[string]*ClientAuthData
	ksid     []byte
}

func NewClientCodec(conn io.ReadWriteCloser) ClientCodec {
	return &clientCodec{
		conn:        conn,
		appID:       fixAppID(),
		networkType: fixNetworkType(defaultDeviceNetworkType),
		netIPFamily: fixNetIPFamily(defaultDeviceNetIPFamily),
		imei:        defaultDeviceIMEI,
		imsi:        defaultDeviceIMSI,
		revision:    defaultClientRevision,
	}
}

func fixAppID() uint32 {
	appID := defaultClientCodecAppIDRelease
	for i := range appID {
		appID[i] ^= defaultClientCodecAppIDMapByte[i%4]
	}
	v, _ := strconv.Atoi(string(appID))
	return uint32(v)
}

func fixNetworkType(v string) uint8 {
	switch v {
	case "Wi-Fi":
		return 0x01
	default:
		return 0x00
	}
}

func fixNetIPFamily(v string) uint8 {
	switch v {
	case "IPv4Only":
		return 0x01
	case "IPv6Only":
		return 0x02
	case "IPv4IPv6":
		return 0x03
	default:
		return 0x00
	}
}

func (c *clientCodec) serializeHead(msg *ClientToServerMessage) ([]byte, error) {
	if msg.Version != 0x0000000a && msg.Version != 0x0000000b {
		return nil, fmt.Errorf("failed to serialize head, version 0x%x", msg.Version)
	}
	buf := bytes.NewBuffer([]byte{})
	buf.EncodeUint32(0x00000000)
	buf.EncodeUint32(msg.Version)
	buf.EncodeUint8(msg.EncryptType)
	switch msg.Version {
	case 0x0000000a:
		buf.EncodeUint32(uint32(len(msg.EncryptKey) + 4))
		buf.EncodeRawBytes(msg.EncryptKey[:])
	case 0x0000000b:
		buf.EncodeUint32(msg.Seq)
	}
	buf.EncodeUint8(0x00)
	buf.EncodeUint32(uint32(len(msg.Username) + 4))
	buf.EncodeRawString(msg.Username)
	return buf.Bytes(), nil
}

func (c *clientCodec) serializeData(msg *ClientToServerMessage) ([]byte, error) {
	if msg.Version != 0x0000000a && msg.Version != 0x0000000b {
		return nil, fmt.Errorf("failed to serialize data, version 0x%x", msg.Version)
	}
	buf := bytes.NewBuffer([]byte{})
	buf.EncodeUint32(0x00000000)
	if msg.Version == 0x0000000a {
		buf.EncodeUint32(msg.Seq)
		buf.EncodeUint32(c.appID)
		buf.EncodeUint32(msg.AppID)
		{
			tmp := make([]byte, 12)
			tmp[0x0] = c.networkType
			tmp[0xa] = c.netIPFamily
			buf.EncodeRawBytes(tmp)
		}
		auth, ok := c.authData[msg.Username]
		if !ok {
			buf.EncodeUint32(4)
		} else {
			buf.EncodeUint32(uint32(len(auth.A2) + 4))
			buf.EncodeRawBytes(auth.A2)
		}
	}
	buf.EncodeUint32(uint32(len(msg.ServiceMethod) + 4))
	buf.EncodeRawString(msg.ServiceMethod)
	buf.EncodeUint32(uint32(len(msg.Cookie) + 4))
	buf.EncodeRawBytes(msg.Cookie)
	if msg.Version == 0x0000000a {
		buf.EncodeUint32(uint32(len(c.imei) + 4))
		buf.EncodeRawString(c.imei)
		buf.EncodeUint32(uint32(len(c.ksid) + 4))
		buf.EncodeRawBytes(c.ksid)
		{
			tmp := "" + "|" + c.imsi + "|A" + c.revision
			buf.EncodeUint32(uint32(len(tmp) + 4))
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

func (c *clientCodec) deserializeHead(buf *bytes.Buffer, msg *ServerToClientMessage) error {
	var err error
	if _, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.Version, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.Version != 0x0000000a && msg.Version != 0x0000000b {
		return fmt.Errorf("failed to deserialize head, version 0x%x", msg.Version)
	}
	if msg.EncryptType, err = buf.DecodeUint8(); err != nil {
		return err
	}
	l, err := buf.DecodeUint32()
	if err != nil {
		return err
	}
	if msg.Username, err = buf.DecodeStringN(uint16(l - 4)); err != nil {
		return err
	}
	return nil
}

func (c *clientCodec) deserializeData(buf *bytes.Buffer, msg *ServerToClientMessage) error {
	if msg.Version != 0x0000000a && msg.Version != 0x0000000b {
		return fmt.Errorf("failed to serialize head, version 0x%x", msg.Version)
	}
	var err error
	if _, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.Seq, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if msg.ReturnCode, err = buf.DecodeUint32(); err != nil {
		return err
	}
	var l uint32
	if l, err = buf.DecodeUint32(); err != nil {
		return err
	}
	if _, err = buf.DecodeBytesN(uint16(l - 4)); err != nil {
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
	if _, err = buf.DecodeUint32(); err != nil {
		return err
	}
	msg.Buffer = buf.Bytes()
	return nil
}

func (c *clientCodec) Encode(msg *ClientToServerMessage) error {
	if !msg.Simple {
		msg.Version = 0x0000000a
	} else {
		msg.Version = 0x0000000b
	}
	data, err := c.serializeData(msg)
	if err != nil {
		return err
	}
	method := strings.ToLower(msg.ServiceMethod)
	if method == "heartbeat.ping" ||
		method == "heartbeat.alive" ||
		method == "client.correcttime" {
		msg.EncryptType = 0x00
	} else {
		_, ok := c.authData[msg.Username]
		cipher := crypto.NewCipher([16]byte{})
		if !ok || len(msg.EncryptKey) == 0 ||
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
			cipher.SetKey(msg.EncryptKey)
			msg.EncryptType = 0x01
		}
		data = cipher.Encrypt(data)
	}
	head, err := c.serializeHead(msg)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(head[0:], uint32(len(head)+len(data)))
	if _, err := c.conn.Write(head); err != nil {
		return err
	}
	if _, err := c.conn.Write(data); err != nil {
		return err
	}
	return nil
}

func (c *clientCodec) Decode(msg *ServerToClientMessage) error {
	v := make([]byte, 4)
	if _, err := c.conn.Read(v); err != nil {
		return err
	}
	l := uint32(v[0])<<24 | uint32(v[1])<<16 | uint32(v[2])<<8 | uint32(v[3])<<0
	v = append(v, make([]byte, l-4)...)
	if _, err := c.conn.Read(v[4:]); err != nil {
		return err
	}
	buf := bytes.NewBuffer(v)
	if err := c.deserializeHead(buf, msg); err != nil {
		return err
	}
	switch msg.EncryptType {
	case 0x00:
	case 0x01:
		auth, ok := c.authData[msg.Username]
		if ok {
			buf = bytes.NewBuffer(crypto.NewCipher(auth.Key).Decrypt(buf.Bytes()))
		}
	case 0x02:
		buf = bytes.NewBuffer(crypto.NewCipher([16]byte{}).Decrypt(buf.Bytes()))
	case 0x03:
	default:
		return fmt.Errorf("failed to decode data, encrypt type 0x%x", msg.EncryptType)
	}
	if err := c.deserializeData(buf, msg); err != nil {
		return err
	}
	return nil
}

func (c *clientCodec) Close() error {
	return c.conn.Close()
}
