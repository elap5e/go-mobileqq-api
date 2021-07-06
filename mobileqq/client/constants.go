package client

var (
	defaultDeviceOSType    = "android"
	defaultDeviceOSVersion = "11"

	defaultDeviceOSBuildBrand = []byte("Xiaomi")
	defaultDeviceOSBuildModel = "Redmi K20"
	defaultDeviceOSSDKVersion = uint32(30)

	defaultDeviceAPNName   = []byte("wifi")
	defaultDeviceSIMOPName = []byte("CMCC")
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
	ServiceMethodAccountGetLoginDevice  = "StatSvc.GetDevLoginInfo"
	ServiceMethodAccountUpdateStatus    = "StatSvc.register"
	ServiceMethodAccountPushLoginNotify = "StatSvc.SvcReqMSFLoginNotify"

	ServiceMethodAuthLogin           = "wtlogin.login"
	ServiceMethodAuthExchange        = "wtlogin.exchange"
	ServiceMethodAuthExchangeAccount = "wtlogin.exchange_emp"
	ServiceMethodAuthUsernameToUin   = "wtlogin.name2uin"
	ServiceMethodAuthRegisterAccount = "wtlogin.trans_emp"

	ServiceMethodConfigPushDomain   = "ConfigPushSvc.PushDomain"
	ServiceMethodConfigPushRequest  = "ConfigPushSvc.PushReq"
	ServiceMethodConfigPushResponse = "ConfigPushSvc.PushResp"

	ServiceMethodMessageDeleteMessage = "MessageSvc.PbDeleteMsg"
	ServiceMethodMessageGetMessage    = "MessageSvc.PbGetMsg"
	ServiceMethodMessageSendMessage   = "MessageSvc.PbSendMsg"
	ServiceMethodMessagePushNotify    = "MessageSvc.PushNotify"
	ServiceMethodMessagePushReaded    = "MessageSvc.PushReaded"

	ServiceMethodOnlinePushMessageSyncC2C   = "OnlinePush.PbC2CMsgSync"
	ServiceMethodOnlinePushMessageSyncGroup = "OnlinePush.PbPushGroupMsg"
	ServiceMethodOnlinePushRequest          = "OnlinePush.ReqPush"
	ServiceMethodOnlinePushResponse         = "OnlinePush.RespPush"
	ServiceMethodOnlinePushSIDTicketExpired = "OnlinePush.SidTicketExpired"
)