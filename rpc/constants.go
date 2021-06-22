package rpc

import (
	"crypto/rand"
)

var (
	defaultDeviceKSID = []byte("")
	defaultDeviceGUID = []byte("%4;7t>;28<fclient.5*6")

	defaultClientRandomKey = func() [16]byte { var v [16]byte; rand.Read(v[:]); return v }()
)

var (
	defaultDeviceOSID   = []byte("")
	defaultDeviceOSType = []byte("android")

	defaultDeviceOSBuildBrand          = []byte("XIAOMI")
	defaultDeviceOSBuildModel          = []byte("MIUI")
	defaultDeviceOSBuildVersionRelease = []byte("10")

	defaultDeviceIMEI = "867345045018141"
	defaultDeviceIMSI = "088976436707562405"

	defaultDeviceAPNName   = []byte("")
	defaultDeviceSIMOPName = []byte("CMCC")

	defaultDeviceNetworkType   = "Wi-Fi"
	defaultDeviceNetworkTypeID = uint16(0x0002)
	defaultDeviceNetIPFamily   = "IPv4IPv6"
	defaultDeviceMACAddress    = []byte("00:05:69:0E:8E:9A")
	defaultDeviceBSSIDAddress  = []byte("00:50:56:9B:EA:5F")
	defaultDeviceSSIDAddress   = []byte("")
)

var (
	defaultAPKID           = []byte("com.tencent.mobileqq")
	defaultAPKSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}
	defaultAPKVersionName  = []byte("8.8.0")
)

var (
	defaultClientCodecAppIDRelease = []byte("7363488;3")
	defaultClientCodecAppIDMapByte = map[int]uint8{0: 2, 1: 0, 2: 1, 3: 3}
)

var (
	defaultClientDstAppID     = uint64(0x0000000000000010)
	defaultClientSubDstAppID  = uint64(0x0000000000000010)
	defaultClientSubAppIDList = []uint64{0x000000005f5e10e2}

	defaultClientMainSigMap = uint32(0x08ff32f2)
	defaultClientMiscBitmap = uint32(0x08f7ff7c)
	defaultClientSubSigMap  = uint32(0x00010400)

	defaultClientBuildTime  = uint64(0x00000000609b85ad)
	defaultClientSDKVersion = "6.0.0.2476"
	defaultClientDomains    = []string{
		"game.qq.com",
		"mail.qq.com",
		"qzone.qq.com",
		"qun.qq.com",
		"openmobile.qq.com",
		"tenpay.com",
		"connect.qq.com",
		"qun.qq.com",
		"qqweb.qq.com",
		"office.qq.com",
		"ti.qq.com",
		"mma.qq.com",
		"docs.qq.com",
		"vip.qq.com",
		"gamecenter.qq.com",
	}

	defaultClientLocaleID = uint32(0x00000804)
	defaultClientRevision = "8.8.0.80d44512"
)
