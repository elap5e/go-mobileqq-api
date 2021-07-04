package tcp

import (
	"encoding/binary"
	"strings"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (c *clientCodec) writeBody(
	msg *codec.ClientToServerMessage,
	buf *bytes.Buffer,
) error {
	buf.WriteUint32(0x00000000) // 0x00000000
	if msg.Version == codec.VersionDefault {
		buf.WriteUint32(msg.Seq)
		buf.WriteUint32(msg.FixID)
		buf.WriteUint32(msg.AppID)
		tmp := make([]byte, 12)
		tmp[0x0] = msg.NetworkType
		tmp[0xa] = msg.NetIPFamily
		buf.Write(tmp)
		buf.WriteUint32Bytes(msg.UserA2)
	}
	buf.WriteUint32String(msg.ServiceMethod)
	buf.WriteUint32Bytes(msg.Cookie)
	if msg.Version == codec.VersionDefault {
		buf.WriteUint32String(msg.IMEI)
		buf.WriteUint32Bytes(msg.KSID)
		buf.WriteUint16String("" + "|" + msg.IMSI + "|A" + msg.Revision)
	}
	buf.WriteUint32Bytes(msg.ReserveField)
	buf.WriteUint32At(uint32(buf.Len()), 0)
	buf.WriteUint32Bytes(msg.Buffer)
	return nil
}

func (c *clientCodec) writeHead(
	msg *codec.ClientToServerMessage,
	buf *bytes.Buffer,
) error {
	buf.WriteUint32(0x00000000) // 0x00000000
	buf.WriteUint32(msg.Version)
	buf.WriteUint8(msg.EncryptType)
	switch msg.Version {
	case codec.VersionDefault:
		if msg.EncryptType == codec.EncryptTypeEncryptByD2Key {
			buf.WriteUint32Bytes(msg.UserD2)
		} else {
			buf.WriteUint32(0x00000004)
		}
	case codec.VersionSimple:
		buf.WriteUint32(msg.Seq)
	}
	buf.WriteUint8(0x00) // 0x00
	buf.WriteUint32String(msg.Username)
	return nil
}

func (c *clientCodec) Write(msg *codec.ClientToServerMessage) error {
	if !msg.Simple {
		msg.Version = codec.VersionDefault
	} else {
		msg.Version = codec.VersionSimple
	}

	bufBody := bufPool.Get()
	defer bufPool.Put(bufBody)
	if err := c.writeBody(msg, bufBody); err != nil {
		return err
	}
	body := bufBody.Bytes()

	method := strings.ToLower(msg.ServiceMethod)
	if method == "heartbeat.ping" ||
		method == "heartbeat.alive" ||
		method == "client.correcttime" {
		msg.EncryptType = codec.EncryptTypeNotNeedEncrypt
	} else {
		cipher := crypto.NewCipher([16]byte{})
		if len(msg.UserD2) == 0 ||
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
			msg.EncryptType = codec.EncryptTypeEncryptByZeros
		} else {
			cipher.SetKey(msg.UserD2Key)
			msg.EncryptType = codec.EncryptTypeEncryptByD2Key
		}
		body = cipher.Encrypt(body)
	}

	bufHead := bufPool.Get()
	defer bufPool.Put(bufHead)
	if err := c.writeHead(msg, bufHead); err != nil {
		return err
	}
	head := bufHead.Bytes()
	binary.BigEndian.PutUint32(head[0:], uint32(len(head)+len(body)))

	if _, err := c.conn.Write(head); err != nil {
		return err
	}
	if _, err := c.conn.Write(body); err != nil {
		return err
	}
	return nil
}
