package rpc

import (
	"encoding/base64"
	"net"
)

var (
	defaultServerECDHPublicKey = []byte{
		0x04, 0xeb, 0xca, 0x94, 0xd7, 0x33, 0xe3, 0x99,
		0xb2, 0xdb, 0x96, 0xea, 0xcd, 0xd3, 0xf6, 0x9a,
		0x8b, 0xb0, 0xf7, 0x42, 0x24, 0xe2, 0xb4, 0x4e,
		0x33, 0x57, 0x81, 0x22, 0x11, 0xd2, 0xe6, 0x2e,
		0xfb, 0xc9, 0x1b, 0xb5, 0x53, 0x09, 0x8e, 0x25,
		0xe3, 0x3a, 0x79, 0x9a, 0xdc, 0x7f, 0x76, 0xfe,
		0xb2, 0x08, 0xda, 0x7c, 0x65, 0x22, 0xcd, 0xb0,
		0x71, 0x9a, 0x30, 0x51, 0x80, 0xcc, 0x54, 0xa8,
		0x2e,
	}
	defaultServerRSAPublicKey, _ = base64.StdEncoding.DecodeString("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuJTW4abQJXeVdAODw1CamZH4QJZChyT08ribet1Gp0wpSabIgyKFZAOxeArcCbknKyBrRY3FFI9HgY1AyItH8DOUe6ajDEb6c+vrgjgeCiOiCVyum4lI5Fmp38iHKH14xap6xGaXcBccdOZNzGT82sPDM2Oc6QYSZpfs8EO7TYT7KSB2gaHz99RQ4A/Lel1Vw0krk+DescN6TgRCaXjSGn268jD7lOO23x5JS1mavsUJtOZpXkK9GqCGSTCTbCwZhI33CpwdQ2EHLhiP5RaXZCio6lksu+d8sKTWU1eEiEb3cQ7nuZXLYH7leeYFoPtbFV4RicIWp0/YG+RP7rLPCwIDAQAB")
)

var (
	defaultDeviceOSType    = "android"
	defaultDeviceOSVersion = "11"

	defaultDeviceOSBuildID    = []byte("RKQ1.200827.002")
	defaultDeviceOSBuildBrand = []byte("Xiaomi")
	defaultDeviceOSBuildModel = "Redmi K20"
	defaultDeviceOSSDKVersion = uint32(30)

	defaultDeviceIMEI = "860308028836598"
	defaultDeviceIMSI = "088906035901507678"

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

const (
	ServiceMethodAuthLogin           = "wtlogin.login"
	ServiceMethodAuthExchange        = "wtlogin.exchange"
	ServiceMethodAuthExchangeAccount = "wtlogin.exchange_emp"
	ServiceMethodAuthUsernameToUin   = "wtlogin.name2uin"
	ServiceMethodAuthRegisterAccount = "wtlogin.trans_emp"

	ServiceMethodAccountUpdateStatus       = "StatSvc.register"
	ServiceMethodAccountGetDeviceLoginInfo = "StatSvc.GetDevLoginInfo"

	ServiceMethodMessageGetMessage = "MessageSvc.PbGetMsg"

	ServiceMethodPushConfigDomain           = "ConfigPushSvc.PushDomain"
	ServiceMethodPushConfigRequest          = "ConfigPushSvc.PushReq"
	ServiceMethodPushConfigResponse         = "ConfigPushSvc.PushResp"
	ServiceMethodPushMessageNotify          = "MessageSvc.PushNotify"
	ServiceMethodPushOnlineSIDTicketExpired = "OnlinePush.SidTicketExpired"
)
