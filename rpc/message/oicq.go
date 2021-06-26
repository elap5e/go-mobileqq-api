package message

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/tlv"
	"github.com/elap5e/go-mobileqq-api/util"
)

type OICQMessage struct {
	Version       uint16
	ServiceMethod uint16
	Uin           uint64
	EncryptMethod uint8
	RandomKey     [16]byte
	KeyVersion    uint16
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
	data, err := marshalOICQMessageData(ctx, msg)
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
	log.Printf("--> [recv] encryptMethod 0x%02x, dump oicq:\n%s", msg.EncryptMethod, hex.Dump(buf.Bytes()))
	if err := unmarshalOICQMessageData(ctx, buf, msg); err != nil {
		return err
	}
	DumpTLVs(ctx, msg.TLVs)
	return nil
}

func DumpTLVs(ctx context.Context, tlvs map[uint16]tlv.TLVCodec, flag ...bool) {
	for i := range tlvs {
		v := tlvs[i].(*tlv.TLV)
		buf, _ := v.GetValue()
		switch i {
		case 0x000a: // message 服务器繁忙，请你稍后再试。
			code, _ := buf.DecodeUint32()
			message, _ := buf.DecodeString()
			log.Printf("--> [recv] dump tlv 0x000a, code 0x%08x, message %s", code, message)
		case 0x0146: // message 登录失败
			code, _ := buf.DecodeUint32()
			title, _ := buf.DecodeString()
			message, _ := buf.DecodeString()
			log.Printf("--> [recv] dump tlv 0x0146, code 0x%08x, title %s, message %s", code, title, message)
		case 0x017d: // type:message,type:message
			type1, _ := buf.DecodeUint16()
			message1, _ := buf.DecodeString()
			type2, _ := buf.DecodeUint16()
			message2, _ := buf.DecodeString()
			log.Printf("--> [recv] dump tlv 0x017d, type1 0x%04x, message1 %s, type2 0x%04x, message2 %s", type1, message1, type2, message2)
		case 0x0104: // session
			log.Printf("--> [recv] dump tlv 0x0104, session %s", string(buf.Bytes()))
		case 0x0178: // countryCode:mobile
			countryCode, _ := buf.DecodeString()
			mobile, _ := buf.DecodeString()
			status, _ := buf.DecodeUint32()
			counts, _ := buf.DecodeUint16()
			timeLimit, _ := buf.DecodeUint16()
			log.Printf("--> [recv] dump tlv 0x0178, country code %s, mobile %s, status 0x%08x, counts 0x%04x, timeLimit 0x%04x", countryCode, mobile, status, counts, timeLimit)
		case 0x0105: // picture
			sign, _ := buf.DecodeBytes()
			data, _ := buf.DecodeBytes()
			log.Printf("--> [recv] dump tlv 0x0105, sign 0x%x, picture length %d", sign, len(data))
		case 0x0165: // picture
			_, _ = buf.DecodeUint32()
			l, _ := buf.DecodeUint8()
			code, _ := buf.DecodeStringN(uint16(l))
			_, _ = buf.DecodeUint16()
			message, _ := buf.DecodeString()
			log.Printf("--> [recv] dump tlv 0x0165, code %s, message %s", code, message)
		case 0x0192: // captcha
			log.Printf("--> [recv] dump tlv 0x0192, url %s", string(buf.Bytes()))

		// TODO: recv 0x0119 sub tlvs
		case 0x0119:
			log.Printf("--> [recv] dump tlv 0x0119(encrypt):\n%s", hex.Dump(buf.Bytes()))
		case 0x0161, 0x0163, 0x0522, 0x0537, 0x0550: // decrypt
			log.Printf("--> [recv] dump tlv 0x%04x(decrypt):\n%s", i, hex.Dump(buf.Bytes()))
		case 0x0102: // clientCodecKeyA8
			log.Printf("--> [recv] dump tlv 0x0102(decrypt), clientCodecKeyA8\n%s", hex.Dump(buf.Bytes()))
		case 0x0106: // _encryptA1=TGTGT
			log.Printf("--> [recv] dump tlv 0x0106(decrypt), tgtgt:\n%s", hex.Dump(buf.Bytes()))
		case 0x010c: // _encryptA1=GTKey
			log.Printf("--> [recv] dump tlv 0x010c(decrypt), gtKey:\n%s", hex.Dump(buf.Bytes()))
		case 0x010a: // clientCodecKeyA2
			log.Printf("--> [recv] dump tlv 0x010a(decrypt), clientCodecKeyA2\n%s", hex.Dump(buf.Bytes()))
		case 0x010b: // clientCodecKeyA5
			log.Printf("--> [recv] dump tlv 0x010b(decrypt), clientCodecKeyA5\n%s", hex.Dump(buf.Bytes()))
		case 0x011d: // sessionTicket
			appId, _ := buf.DecodeUint32()
			stKey, _ := buf.DecodeBytesN(0x0016)
			st, _ := buf.DecodeBytes()
			log.Printf("--> [recv] dump tlv 0x011d(decrypt), appID 0x%08x", appId)
			log.Printf("--> [recv] dump tlv 0x011d(decrypt), sessionTicketKey\n%s", hex.Dump(stKey))
			log.Printf("--> [recv] dump tlv 0x011d(decrypt), sessionTicket\n%s", hex.Dump(st))
		case 0x011f: // tk_pri
			chgt, _ := buf.DecodeUint32()
			tk, _ := buf.DecodeUint32()
			log.Printf("--> [recv] dump tlv 0x011f(decrypt), change time %s, tk_pri 0x%08x", time.Unix(int64(chgt), 0), tk)
		case 0x0130: // current server time, ip address
			_, _ = buf.DecodeUint16()
			sct, _ := buf.DecodeUint32()
			util.SetServerCurrentTime(sct)
			ip, _ := buf.DecodeBytesN(0x0004)
			log.Printf("--> [recv] dump tlv 0x010a(decrypt), current server time %s, ip %s", time.Unix(int64(sct), 0), net.IPv4(ip[0], ip[1], ip[2], ip[3]))
		case 0x0134: // wtSessionTicketKey
			log.Printf("--> [recv] dump tlv 0x0134(decrypt), wtSessionTicketKey:\n%s", hex.Dump(buf.Bytes()))
		case 0x0143: // clientCodecKeyD2
			log.Printf("--> [recv] dump tlv 0x0143(decrypt), clientCodecKeyD2:\n%s", hex.Dump(buf.Bytes()))
		case 0x016a: // _noPictureSignature
			log.Printf("--> [recv] dump tlv 0x016a(decrypt), _noPictureSignature:\n%s", hex.Dump(buf.Bytes()))
		case 0x016d: // _superKey
			log.Printf("--> [recv] dump tlv 0x016d(decrypt), _superKey %s", string(buf.Bytes()))
		case 0x0512: // domain tickets
			l, _ := buf.DecodeUint16()
			kv := map[string]string{}
			for i := 0; i < int(l); i++ {
				key, _ := buf.DecodeString()
				kv[key], _ = buf.DecodeString()
				_, _ = buf.DecodeUint16()
			}
			log.Printf("--> [recv] dump tlv 0x0512(decrypt), domain tickets:")
			for key := range kv {
				fmt.Printf("domain %17s, ticket %s\n", key, kv[key])
			}
		case 0x0528: // _loginResultField1
			log.Printf("--> [recv] dump tlv 0x0528(decrypt), _loginResultField1 %s", string(buf.Bytes()))

		case 0x0005: // _psKey
		case 0x0103: // userSessionTicketWeb
		case 0x010d: // clientCodecKeyA2Key
		case 0x010e: // userSessionTicketKey
		case 0x0114: // userSessionTicket
		case 0x011c: // _lsKey
		case 0x0120: // _sKey
		case 0x0121: // _userSig64
		case 0x0125: // _openKey, _openID
		case 0x0132: // _accessToken, _openID
		case 0x0133: // wtSessionTicket
		case 0x0136: // _vKey
		case 0x0164: // _sid
		case 0x0171: // _aqSig
		case 0x0199: // _payToken, _openID
		case 0x0200: // _pf, _pfKey
		case 0x0203: // _da2
		case 0x0305: // clientCodecKeyD2Key
		case 0x0322: // _deviceToken
		case 0x0403: // randseed
			// bArr19[2] G
			// bArr19[3] DPWD
		case 0x0530: // _loginResultField2
			log.Printf("--> [recv] dump tlv 0x0530, _loginResultField2 %s", string(buf.Bytes()))

		case 0x0186: // pwdFlag
		case 0x0108: // ksid
		case 0x0118: // mainDisplayName
		case 0x011a: // face, age, gender, nick
		case 0x0167: // reserveUinInfo
		case 0x0138: // {a2, lsKey, sKey, vKey, a8, stWeb, d2, sid}ChangeTime

		case 0x0179: // 2 bytes
			log.Printf("--> [recv] dump tlv 0x%04x, 2 bytes:\n%s", i, hex.Dump(buf.Bytes()))
		// case 0x0402, 0x0403: // 4 bytes
		// 	log.Printf("--> [recv] dump tlv 0x%04x, 8 byte:\n%s", i, hex.Dump(buf.Bytes()))
		case 0x017e, 0x0204: // string
			log.Printf("--> [recv] dump tlv 0x%04x, raw %s", i, string(buf.Bytes()))
		default:
			log.Printf("--> [recv] dump tlv 0x%04x:\n%s", i, hex.Dump(buf.Bytes()))
		}
		buf.Seek(0)
	}
}
