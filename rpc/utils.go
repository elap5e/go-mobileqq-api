package rpc

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

var (
	deviceGUID          = md5.Sum(append(defaultDeviceOSBuildID, defaultDeviceMACAddress...)) // []byte("%4;7t>;28<fclient.5*6")
	deviceGUIDFlag      = uint32((1 << 24 & 0xFF000000) | (0 << 8 & 0xFF00))
	deviceIsGUIDFileNil = false
	deviceIsGUIDGenSucc = true
	deviceIsGUIDChanged = false

	clientVerifyMethod = uint8(0x82) // 0x00, 0x82
)

var (
	clientPackageName  = []byte("com.tencent.mobileqq")
	clientVersionName  = []byte("8.8.3")
	clientRevision     = "8.8.3.b2791edc"
	clientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}

	clientAppID      = uint32(0x20030cb2)
	clientBuildTime  = uint64(0x00000000609b85ad)
	clientSDKVersion = "6.0.0.2476"
	clientSSOVersion = uint32(0x00000011)

	clientCodecAppIDDebug   = []byte("736350642")
	clientCodecAppIDRelease = []byte("736350642")

	clientImageType  = uint8(0x01)
	clientMiscBitmap = uint32(0x08f7ff7c)
)

func SetClientForAndroid() {
	clientPackageName = []byte("com.tencent.mobileqq")
	clientVersionName = []byte("8.8.3")
	clientRevision = "8.8.3.b2791edc"
	clientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}

	clientAppID = uint32(0x20030cb2)
	clientBuildTime = uint64(0x00000000609b85ad)
	clientSDKVersion = "6.0.0.2476"
	clientSSOVersion = uint32(0x00000011)

	clientCodecAppIDDebug = []byte("736350642")
	clientCodecAppIDRelease = []byte("736350642")

	clientMiscBitmap = uint32(0x08f7ff7c)
	tlv.SetSSOVersion(clientSSOVersion)
}

func SetClientLiteForAndroid() {
	clientPackageName = []byte("com.tencent.qqlite")
	clientVersionName = []byte("4.0.2")
	clientRevision = "4.0.2.9b6340cd"
	clientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}

	clientAppID = uint32(0x200300f0)
	clientBuildTime = uint64(0x0000000060409d2d)
	clientSDKVersion = "6.0.0.2356"
	clientSSOVersion = uint32(0x00000005)

	clientCodecAppIDDebug = []byte("736360370")
	clientCodecAppIDRelease = []byte("736347652")

	clientMiscBitmap = uint32(0x00f7ff7c)
	tlv.SetSSOVersion(clientSSOVersion)
}

func SetClientForAndroidTablet() {
	clientPackageName = []byte("com.tencent.minihd.qq")
	clientVersionName = []byte("5.9.2")
	clientRevision = "5.9.2.3baec0"
	clientSignatureMD5 = [16]byte{0xaa, 0x39, 0x78, 0xf4, 0x1f, 0xd9, 0x6f, 0xf9, 0x91, 0x4a, 0x66, 0x9e, 0x18, 0x64, 0x74, 0xc7}

	clientAppID = uint32(0x2002fdd5)
	clientBuildTime = uint64(0x000000005f1e8730)
	clientSDKVersion = "6.0.0.2433"
	clientSSOVersion = uint32(0x0000000c)

	clientCodecAppIDDebug = []byte("73636270;")
	clientCodecAppIDRelease = []byte("736346857")

	clientMiscBitmap = uint32(0x08f7ff7c)
	tlv.SetSSOVersion(clientSSOVersion)
}

func SetClientForAndroidWatch() {
	panic("not implement")
}

func SetClientForiOS() {
	panic("not implement")
}

func SetClientLiteForiOS() {
	panic("not implement")
}

func SetClientForiPadOS() {
	panic("not implement")
}

func SetClientForwatchOS() {
	panic("not implement")
}

func SetClientForWindows() {
	panic("not implement")
}

func SetClientFormacOS() {
	panic("not implement")
}

func SetClientForLinux() {
	panic("not implement")
}

func ParseUserSignature(ctx context.Context, username string, tlvs map[uint16]tlv.TLVCodec) *UserSignature {
	token := []byte{}
	if v, ok := tlvs[0x0322]; ok {
		token = v.(*tlv.TLV).MustGetValue().Bytes()
	}

	domains := map[string]string{}
	{
		buf := bytes.NewBuffer(tlvs[0x0512].(*tlv.TLV).MustGetValue().Bytes())
		l, _ := buf.DecodeUint16()
		for i := 0; i < int(l); i++ {
			key, _ := buf.DecodeString()
			domains[key], _ = buf.DecodeString()
			_, _ = buf.DecodeUint16()
		}
	}

	chgt := map[uint16]uint32{}
	{
		buf := bytes.NewBuffer(tlvs[0x0138].(*tlv.TLV).MustGetValue().Bytes())
		l, _ := buf.DecodeUint32()
		for i := 0; i < int(l); i++ {
			key, _ := buf.DecodeUint16()
			chgt[key], _ = buf.DecodeUint32()
			_, _ = buf.DecodeUint32()
		}
	}

	tickets := map[string]Ticket{}
	{
		if v, ok := tlvs[0x010a]; ok {
			tickets["A2"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: tlvs[0x010d].(*tlv.TLV).MustGetValue().Bytes(),
				Iss: time.Now().Unix(),
				Exp: time.Now().Unix() + int64(chgt[0x010a]),
			}
		}
		if v, ok := tlvs[0x010b]; ok {
			tickets["A5"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: -1,
			}
		}
		if v, ok := tlvs[0x0102]; ok {
			tickets["A8"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: -1,
			}
		}
		if v, ok := tlvs[0x0143]; ok {
			tickets["D2"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: tlvs[0x0305].(*tlv.TLV).MustGetValue().Bytes(),
				Iss: time.Now().Unix(),
				Exp: time.Now().Unix() + int64(chgt[0x0143]),
			}
		}
		if v, ok := tlvs[0x011c]; ok {
			tickets["LSKey"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: time.Now().Unix() + int64(chgt[0x011c]),
			}
		}
		if v, ok := tlvs[0x0120]; ok {
			tickets["SKey"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: time.Now().Unix() + int64(chgt[0x0120]),
			}
		}
		if v, ok := tlvs[0x0164]; ok {
			tickets["Sig64"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: -1,
			}
		}
		if v, ok := tlvs[0x0164]; ok {
			tickets["SID"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: time.Now().Unix() + int64(chgt[0x0164]),
			}
		}
		if v, ok := tlvs[0x0114]; ok {
			tickets["ST"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: tlvs[0x010e].(*tlv.TLV).MustGetValue().Bytes(),
				Iss: time.Now().Unix(),
				Exp: -1,
			}
		}
		if v, ok := tlvs[0x0103]; ok {
			tickets["STWeb"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: time.Now().Unix() + int64(chgt[0x0103]),
			}
		}
		if v, ok := tlvs[0x0136]; ok {
			tickets["VKey"] = Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: time.Now().Unix() + int64(chgt[0x0136]),
			}
		}
		// 0x00004000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: -1,
		// }, // ??? OpenKey 0x0125
		// 0x00008000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: -1,
		// }, // ??? AccessToken 0x0132
		// 0x00100000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: -1,
		// }, // ??? SuperKey 0x016d
		// 0x00200000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: -1,
		// }, // ??? AQSig 0x0171
		// 0x00800000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: -1,
		// }, // ??? PayToken 0x0199
		// 0x01000000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: -1,
		// }, // ??? PF 0x0200
		// 0x02000000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: -1,
		// }, // ??? DA2 0x0203
	}

	return &UserSignature{
		Username:    username,
		DeviceToken: token,
		Domains:     domains,
		Tickets:     tickets,
	}
}

func SetUserSignature(sig *UserSignature) {
	if sig != nil {
		data, _ := ioutil.ReadFile("tmp/user_signature.json")
		tsig := new(UserSignature)
		_ = json.Unmarshal(data, tsig)
		if len(sig.DeviceToken) != 0 {
			tsig.DeviceToken = sig.DeviceToken
		}
		for k, v := range sig.Domains {
			tsig.Domains[k] = v
		}
		for k, v := range sig.Tickets {
			tsig.Tickets[k] = v
		}
		file, _ := json.MarshalIndent(tsig, "", "    ")
		_ = ioutil.WriteFile("tmp/user_signature.json", file, 0600)
	}
}

func GetUserSignature() *UserSignature {
	file, _ := ioutil.ReadFile("tmp/user_signature.json")
	sig := new(UserSignature)
	_ = json.Unmarshal(file, sig)
	return sig
}
