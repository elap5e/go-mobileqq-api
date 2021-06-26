package tlv

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/elap5e/go-mobileqq-api/util"
)

func DumpTLVs(ctx context.Context, tlvs map[uint16]TLVCodec, flag ...bool) {
	for i := range tlvs {
		v := tlvs[i].(*TLV)
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

		case 0x0104: // session
			log.Printf("--> [recv] dump tlv 0x0104, session %s", string(buf.Bytes()))
		case 0x0402: // 4 bytes
			log.Printf("--> [recv] dump tlv 0x%04x, 8 byte:\n%s", i, hex.Dump(buf.Bytes()))

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
		case 0x017e, 0x0204: // string
			log.Printf("--> [recv] dump tlv 0x%04x, raw %s", i, string(buf.Bytes()))
		default:
			log.Printf("--> [recv] dump tlv 0x%04x:\n%s", i, hex.Dump(buf.Bytes()))
		}
		buf.Seek(0)
	}
}
