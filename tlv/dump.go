package tlv

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/util"
)

func DumpTLVs(ctx context.Context, tlvs map[uint16]TLVCodec, flag ...bool) {
	for i := range tlvs {
		v := tlvs[i].(*TLV)
		buf, _ := v.GetValue()
		switch i {
		// error message
		case 0x000a:
			ver, _ := buf.DecodeUint16()
			code, _ := buf.DecodeUint16()
			msg, _ := buf.DecodeString()
			log.Trace().Msgf("--> [recv] tlv 0x000a, ver 0x%04x, code 0x%04x, message %s", ver, code, msg)
		case 0x0146:
			ver, _ := buf.DecodeUint16()
			code, _ := buf.DecodeUint16()
			title, _ := buf.DecodeString()
			msg, _ := buf.DecodeString()
			typ, _ := buf.DecodeUint16()
			link, _ := buf.DecodeString()
			log.Trace().Msgf("--> [recv] tlv 0x0146, ver 0x%04x, code 0x%04x, title %s, message %s, type 0x%04x, link %s", ver, code, title, msg, typ, link)
		case 0x017d:
			code, _ := buf.DecodeUint16()
			msg, _ := buf.DecodeString()
			typ, _ := buf.DecodeUint16()
			link, _ := buf.DecodeString()
			log.Trace().Msgf("--> [recv] tlv 0x017d, code 0x%04x, message %s, type 0x%04x, link %s", code, msg, typ, link)

		case 0x0508: // ???
			log.Trace().Msgf("--> [recv] tlv 0x0508, ???")

		case 0x0105: // picture
			sign, _ := buf.DecodeBytes()
			data, _ := buf.DecodeBytes()
			log.Trace().Msgf("--> [recv] tlv 0x0105, sign 0x%x, picture length %d", sign, len(data))
		case 0x0165: // picture
			_, _ = buf.DecodeUint32()
			l, _ := buf.DecodeUint8()
			code, _ := buf.DecodeStringN(uint16(l))
			_, _ = buf.DecodeUint16()
			message, _ := buf.DecodeString()
			log.Trace().Msgf("--> [recv] tlv 0x0165, code %s, message %s", code, message)
		case 0x0178: // countryCode:mobile
			countryCode, _ := buf.DecodeString()
			mobile, _ := buf.DecodeString()
			status, _ := buf.DecodeUint32()
			counts, _ := buf.DecodeUint16()
			timeLimit, _ := buf.DecodeUint16()
			log.Trace().Msgf("--> [recv] tlv 0x0178, country code %s, mobile %s, status 0x%08x, counts 0x%04x, timeLimit 0x%04x", countryCode, mobile, status, counts, timeLimit)
		case 0x0192: // captcha
			log.Trace().Msgf("--> [recv] tlv 0x0192, url %s", string(buf.Bytes()))

		case 0x0104: // session
			log.Trace().Msgf("--> [recv] tlv 0x0104, session %s", string(buf.Bytes()))
		case 0x0402: // 8 bytes
			log.Trace().Msgf("--> [recv] tlv 0x0402, 8 byte")

		// TODO: recv 0x0119 sub tlvs
		case 0x0119:
			log.Trace().Msgf("--> [recv] tlv 0x0119(encrypt)")
		case 0x0161, 0x0163, 0x0522, 0x0550: // decrypt
			log.Trace().Msgf("--> [recv] tlv 0x%04x(decrypt)", i)
		case 0x0102: // userA8
			log.Trace().Msgf("--> [recv] tlv 0x0102(decrypt), userA8")
		case 0x0103: // userSTWeb
			log.Trace().Msgf("--> [recv] tlv 0x0103(decrypt), userSTWeb")
		case 0x0106: // encryptA1=TGTGT
			log.Trace().Msgf("--> [recv] tlv 0x0106(decrypt), userA1")
		case 0x0108: // ksid
			log.Trace().Msgf("--> [recv] tlv 0x0108(decrypt), ksid")
		case 0x010a: // userA2=TGT
			log.Trace().Msgf("--> [recv] tlv 0x010a(decrypt), userA2")
		case 0x010c: // encryptA1=GTKey
			log.Trace().Msgf("--> [recv] tlv 0x010c(decrypt), userA1Key")
		case 0x010d: // userA2Key
			log.Trace().Msgf("--> [recv] tlv 0x010d(decrypt), userA2Key")
		case 0x010e: // userSTKey
			log.Trace().Msgf("--> [recv] tlv 0x010e(decrypt), userSTKey")
		case 0x010b: // userA5
			log.Trace().Msgf("--> [recv] tlv 0x010b(decrypt), userA5")
		case 0x0114: // userST
			log.Trace().Msgf("--> [recv] tlv 0x0114(decrypt), userST")
		case 0x0118: // mainDisplayName
			log.Trace().Msgf("--> [recv] tlv 0x0118(decrypt), mainDisplayName")
		case 0x011a: // face, age, gender, nick
			log.Trace().Msgf("--> [recv] tlv 0x011a(decrypt), face, age, gender, nick")
		case 0x011c: // userLSKey
			log.Trace().Msgf("--> [recv] tlv 0x011c(decrypt), userLSKey")
		case 0x011d: // sessionTicket
			appId, _ := buf.DecodeUint32()
			stKey, _ := buf.DecodeBytesN(0x0010)
			st, _ := buf.DecodeBytes()
			log.Trace().Msgf("--> [recv] tlv 0x011d(decrypt), appID 0x%08x", appId)
			log.Trace().Msgf("--> [recv] tlv 0x011d(decrypt), sessionTicketKey\n%s", hex.Dump(stKey))
			log.Trace().Msgf("--> [recv] tlv 0x011d(decrypt), sessionTicket\n%s", hex.Dump(st))
		case 0x011f: // tk_pri
			chgt, _ := buf.DecodeUint32()
			tk, _ := buf.DecodeUint32()
			log.Trace().Msgf("--> [recv] tlv 0x011f(decrypt), change time %s, tk_pri 0x%08x", time.Unix(int64(util.GetServerTime()+chgt), 0), tk)
		case 0x0120: // userSKey
			log.Trace().Msgf("--> [recv] tlv 0x0120(decrypt), userSKey %s", string(buf.Bytes()))
		case 0x0130: // current server time, ip address
			_, _ = buf.DecodeUint16()
			svrt, _ := buf.DecodeUint32()
			util.SetServerTime(svrt)
			ip, _ := buf.DecodeBytesN(0x0004)
			log.Trace().Msgf("--> [recv] tlv 0x0130(decrypt), current server time %s, ip %s", time.Unix(int64(svrt), 0), net.IPv4(ip[0], ip[1], ip[2], ip[3]))
		case 0x0133: // wtSessionTicket
			log.Trace().Msgf("--> [recv] tlv 0x0133(decrypt), wtSessionTicket")
		case 0x0134: // wtSessionTicketKey
			log.Trace().Msgf("--> [recv] tlv 0x0134(decrypt), wtSessionTicketKey")
		case 0x0138: // {a2, lsKey, sKey, vKey, a8, stWeb, d2, sid}ChangeTime
			l, _ := buf.DecodeUint32()
			chgt := map[uint16]uint32{}
			for i := 0; i < int(l); i++ {
				key, _ := buf.DecodeUint16()
				chgt[key], _ = buf.DecodeUint32()
				_, _ = buf.DecodeUint32()
			}
			log.Trace().Msgf("--> [recv] tlv 0x0138(decrypt), {a2, lsKey, sKey, vKey, a8, stWeb, d2, sid}ChangeTime")
			log.Trace().Msgf("--> [recv] tlv 0x0138(decrypt),    a2 %s", time.Unix(int64(util.GetServerTime()+chgt[0x010a]), 0))
			log.Trace().Msgf("--> [recv] tlv 0x0138(decrypt), lsKey %s", time.Unix(int64(util.GetServerTime()+chgt[0x011c]), 0))
			log.Trace().Msgf("--> [recv] tlv 0x0138(decrypt),  sKey %s", time.Unix(int64(util.GetServerTime()+chgt[0x0120]), 0))
			log.Trace().Msgf("--> [recv] tlv 0x0138(decrypt),  vKey %s", time.Unix(int64(util.GetServerTime()+chgt[0x0136]), 0))
			log.Trace().Msgf("--> [recv] tlv 0x0138(decrypt),    a8 %s", time.Unix(int64(util.GetServerTime()+chgt[0x0102]), 0))
			log.Trace().Msgf("--> [recv] tlv 0x0138(decrypt), stWeb %s", time.Unix(int64(util.GetServerTime()+chgt[0x0103]), 0))
			log.Trace().Msgf("--> [recv] tlv 0x0138(decrypt),    d2 %s", time.Unix(int64(util.GetServerTime()+chgt[0x0143]), 0))
			log.Trace().Msgf("--> [recv] tlv 0x0138(decrypt),   sid %s", time.Unix(int64(util.GetServerTime()+chgt[0x0164]), 0))
		case 0x0143: // userD2
			log.Trace().Msgf("--> [recv] tlv 0x0143(decrypt), userD2")
		case 0x016a: // noPictureSignature
			log.Trace().Msgf("--> [recv] tlv 0x016a(decrypt), noPictureSignature")
		case 0x016d: // userSuperKey
			log.Trace().Msgf("--> [recv] tlv 0x016d(decrypt), userSuperKey %s", string(buf.Bytes()))
		case 0x0305: // userD2Key
			log.Trace().Msgf("--> [recv] tlv 0x0305(decrypt), userD2Key")
		case 0x0537: // loginExtraData
			log.Trace().Msgf("--> [recv] tlv 0x0537(decrypt), loginExtraData\n%s", hex.Dump(buf.Bytes()[2:]))
		case 0x0512: // domain tickets
			l, _ := buf.DecodeUint16()
			kv := map[string]string{}
			for i := 0; i < int(l); i++ {
				key, _ := buf.DecodeString()
				kv[key], _ = buf.DecodeString()
				_, _ = buf.DecodeUint16()
			}
			tickets := ""
			for key := range kv {
				tickets += fmt.Sprintf("\ndomain %17s, ticket %s", key, kv[key])
			}
			log.Trace().Msgf("--> [recv] tlv 0x0512(decrypt), domain tickets: %s", tickets)
		case 0x0528: // _loginResultField1
			log.Trace().Msgf("--> [recv] tlv 0x0528(decrypt), _loginResultField1 %s", string(buf.Bytes()))

		case 0x0005: // _psKey
		case 0x0121: // userSig64
		case 0x0125: // userOpenKey, userOpenID
		case 0x0132: // userAccessToken, userOpenID
		case 0x0136: // userVKey
		case 0x0164: // userSID
		case 0x0171: // userAQSig
		case 0x0199: // userPayToken, userOpenID
		case 0x0200: // userPF, userPFKey
		case 0x0203: // userDA2
		case 0x0322: // userDeviceToken
		case 0x0403: // randseed
			// bArr19[2] G t401
			// bArr19[3] DPWD
		case 0x0530: // _loginResultField2
			log.Trace().Msgf("--> [recv] tlv 0x0530, _loginResultField2 %s", string(buf.Bytes()))

		case 0x0186: // pwdFlag
		case 0x0167: // reserveUinInfo

		case 0x0179: // 2 bytes
			log.Trace().Msgf("--> [recv] tlv 0x0179, 2 bytes")
		case 0x017e, 0x0204: // string
			log.Trace().Msgf("--> [recv] tlv 0x%04x, raw %s", i, string(buf.Bytes()))
		}
		buf.Seek(0)
		log.Trace().Msgf("--> [recv] dump tlv 0x%04x:\n%s", i, hex.Dump(buf.Bytes()))
	}
}
