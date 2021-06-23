package rpc

import (
	"math/rand"
	"net"
)

var (
	defaultDeviceKSID          = []byte("")
	defaultDeviceGUID          = []byte("%4;7t>;28<fclient.5*6")
	defaultDeviceGUIDFlag      = uint32((1 << 24 & 0xFF000000) | (0 << 8 & 0xFF00))
	defaultDeviceIsGUIDFileNil = false
	defaultDeviceIsGUIDGenSucc = true
	defaultDeviceIsGUIDChanged = false

	defaultClientVerifyMethod = uint8(0x82)

	defaultDeviceDPWD = func(n int) []byte {
		b := make([]byte, n)
		for i := range b {
			b[i] = byte(0x41 + rand.Intn(1)*0x20 + rand.Intn(26))
		}
		return b
	}(16)

	defaultClientRandomKey = func() [16]byte { var v [16]byte; rand.Read(v[:]); return v }()
)

var (
	defaultDeviceOSID      = []byte("")
	defaultDeviceOSType    = []byte("android")
	defaultDeviceOSVersion = []byte("10")

	defaultDeviceOSBuildBrand = []byte("XIAOMI")
	defaultDeviceOSBuildModel = []byte("MIUI")

	defaultDeviceIMEI = "867345045018141"
	defaultDeviceIMSI = "088976436707562405"

	defaultDeviceAPNName   = []byte("")
	defaultDeviceSIMOPName = []byte("CMCC")

	defaultDeviceNetworkType   = "Wi-Fi"
	defaultDeviceNetworkTypeID = uint16(0x0002)
	defaultDeviceNetIPFamily   = "IPv4IPv6"
	defaultDeviceIPv4Address   = net.IPv4(192, 168, 0, 100)
	defaultDeviceMACAddress    = []byte("00:05:69:0E:8E:9A")
	defaultDeviceBSSIDAddress  = []byte("00:50:56:9B:EA:5F")
	defaultDeviceSSIDAddress   = []byte("")
)

var (
	defaultClientPackageName  = []byte("com.tencent.mobileqq")
	defaultClientRevision     = "8.8.3.b2791edc"
	defaultClientVersionName  = []byte("8.8.3")
	defaultClientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}
)

var (
	defaultClientMainSigMap = uint32(0x00ff32f2)
	defaultClientMiscBitmap = uint32(0x08f7ff7c)
	defaultClientSubSigMap  = uint32(0x00010400)
	defaultClientOpenAppID  = uint64(0x000000002a9e5427)
)

var (
	defaultClientCodecAppIDRelease = []byte("736350642")
	defaultClientCodecAppIDMapByte = map[int]uint8{0: 2, 1: 0, 2: 1, 3: 3}
)

var (
	defaultClientBuildTime  = uint64(0x00000000609b85ad)
	defaultClientSDKVersion = "6.0.0.2476"
	defaultClientSSOVersion = uint32(0x00000011)
)

var (
	defaultClientSMSAppID     = uint64(0x0000000000000009)
	defaultClientDstAppID     = uint64(0x0000000000000010)
	defaultClientOpenSDKID    = uint64(0x000000005f5e1604)
	defaultClientSubAppIDList = []uint64{0x000000005f5e10e2}
)

var (
	defaultClientLocaleID = uint32(0x00000804)
	defaultClientDomains  = []string{
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
)
