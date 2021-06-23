package rpc

import (
	"crypto/md5"
	"math/rand"
	"net"
)

var (
	defaultDeviceKSID          = []byte("")
	defaultDeviceGUID          = md5.Sum(append(defaultDeviceOSBuildID, defaultDeviceMACAddress...)) // []byte("%4;7t>;28<fclient.5*6")
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
	defaultDeviceOSBuildID = []byte("8b5c5fe9-a7b8-4f23-bc84-5e6730254da5")
	defaultDeviceOSType    = []byte("android")
	defaultDeviceOSVersion = []byte("10")

	defaultDeviceOSBuildBrand = []byte("be257ffb-865b-4e88-a4ea-5f9ef8617ad0")
	defaultDeviceOSBuildModel = []byte("7c712868-0cad-44d1-9029-dadd6b0f9c21")

	defaultDeviceIMEI = "c5413b37-eade-4ccd-976f-76a659d51684"
	defaultDeviceIMSI = "4c9762cb-d6ae-4ebf-965d-0fc556dc9bc9"

	defaultDeviceAPNName   = []byte("wifi")
	defaultDeviceSIMOPName = []byte("ffb4190b-00f5-4c7b-b3a0-490f58401a5d")

	defaultDeviceNetworkType   = "Wi-Fi"
	defaultDeviceNetworkTypeID = uint16(0x0002)
	defaultDeviceNetIPFamily   = "IPv4IPv6"
	defaultDeviceIPv4Address   = net.IPv4(192, 168, 0, 100)
	defaultDeviceMACAddress    = []byte("00:50:56:C0:00:08")
	defaultDeviceBSSIDAddress  = []byte("00:50:56:C0:00:09")
	defaultDeviceSSIDAddress   = []byte("021dbef5-ca98-4774-ba37-59a844408e39")
)

var (
	defaultClientPackageName  = []byte("com.tencent.mobileqq")
	defaultClientRevision     = "8.8.3.b2791edc"
	defaultClientVersionName  = []byte("8.8.3")
	defaultClientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}
)

var (
	defaultClientAppID      = uint32(0x20030cb2)
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
