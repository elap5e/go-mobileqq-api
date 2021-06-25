package rpc

import (
	"net"
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
	defaultClientMainSigMap        = uint32(0x00ff32f2)
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
