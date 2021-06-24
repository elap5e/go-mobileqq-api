package rpc

import (
	"crypto/md5"
	"math/rand"
	"net"

	"github.com/elap5e/go-mobileqq-api/tlv"
)

func SetClientForAndroidPhone() {
	defaultClientPackageName = []byte("com.tencent.mobileqq")
	defaultClientVersionName = []byte("8.5.0")
	defaultClientRevision = "8.2.7.27f6ea96"
	defaultClientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}
	defaultClientAppID = uint32(0x2002fcf2)
	defaultClientBuildTime = uint64(0x000000005fd36704)
	defaultClientSDKVersion = "6.0.0.2454"
	defaultClientSSOVersion = uint32(0x0000000f)

	defaultClientCodecAppIDDebug = []byte("73634660:")
	defaultClientCodecAppIDRelease = []byte("73634660:")
	tlv.SetSSOVersion(defaultClientSSOVersion)
}

func SetClientForAndroidPad() {
	defaultClientPackageName = []byte("com.tencent.minihd.qq")
	defaultClientVersionName = []byte("5.9.2")
	defaultClientRevision = "5.9.2.3baec0"
	defaultClientSignatureMD5 = [16]byte{0xaa, 0x39, 0x78, 0xf4, 0x1f, 0xd9, 0x6f, 0xf9, 0x91, 0x4a, 0x66, 0x9e, 0x18, 0x64, 0x74, 0xc7}

	defaultClientAppID = uint32(0x2002fdd5)
	defaultClientBuildTime = uint64(0x00000000609b85ad)
	defaultClientSDKVersion = "6.0.0.2433"
	defaultClientSSOVersion = uint32(0x0000000c)

	defaultClientCodecAppIDDebug = []byte("73636270;")
	defaultClientCodecAppIDRelease = []byte("736346857")
	tlv.SetSSOVersion(defaultClientSSOVersion)
}

func SetClientForiPhone() {}

func SetClientForiPad() {}

var (
	defaultDeviceKSID          = []byte("")
	defaultDeviceGUID          = md5.Sum(append(defaultDeviceOSBuildID, defaultDeviceMACAddress...)) // []byte("%4;7t>;28<fclient.5*6")
	defaultDeviceGUIDFlag      = uint32((1 << 24 & 0xFF000000) | (0 << 8 & 0xFF00))
	defaultDeviceIsGUIDFileNil = false
	defaultDeviceIsGUIDGenSucc = true
	defaultDeviceIsGUIDChanged = false

	defaultClientVerifyMethod = uint8(0x82) // 0x00, 0x82

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
	defaultDeviceOSType    = []byte("android")
	defaultDeviceOSVersion = []byte("11")

	defaultDeviceOSBuildID    = []byte("RKQ1.200826.002")
	defaultDeviceOSBuildBrand = []byte("Xiaomi")
	defaultDeviceOSBuildModel = []byte("Redmi K20")

	defaultDeviceIMEI = "860308028836597"
	defaultDeviceIMSI = "088906035901507677"

	defaultDeviceAPNName   = []byte("wifi")
	defaultDeviceSIMOPName = []byte("CMCC")

	defaultDeviceNetworkType   = "Wi-Fi"
	defaultDeviceNetworkTypeID = uint16(0x0002)
	defaultDeviceNetIPFamily   = "IPv4IPv6"
	defaultDeviceIPv4Address   = net.IPv4(192, 168, 0, 100)
	defaultDeviceMACAddress    = []byte("00:50:56:C0:00:08")
	defaultDeviceBSSIDAddress  = []byte("00:50:56:C0:00:09")
	defaultDeviceSSIDAddress   = []byte("unknown")
)

var (
	defaultClientPackageName  = []byte("com.tencent.mobileqq")
	defaultClientVersionName  = []byte("8.8.3")
	defaultClientRevision     = "8.8.3.b2791edc"
	defaultClientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}
)

var (
	defaultClientAppID      = uint32(0x20030cb2)
	defaultClientBuildTime  = uint64(0x00000000609b85ad)
	defaultClientSDKVersion = "6.0.0.2476"
	defaultClientSSOVersion = uint32(0x00000011)

	defaultClientCodecAppIDDebug   = []byte("736350642")
	defaultClientCodecAppIDRelease = []byte("736350642")
)

var (
	defaultClientMainSigMap        = uint32(0x00ff32f2)
	defaultClientMiscBitmap        = uint32(0x08f7ff7c)
	defaultClientSubSigMap         = uint32(0x00010400)
	defaultClientOpenAppID         = uint64(0x000000002a9e5427)
	defaultClientCodecAppIDMapByte = map[int]uint8{0: 2, 1: 0, 2: 1, 3: 3}
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
