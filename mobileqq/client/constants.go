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

const (
	ServiceMethodAccountGetLoginDevice  = "StatSvc.GetDevLoginInfo"
	ServiceMethodAccountUpdateStatus    = "StatSvc.register"
	ServiceMethodAccountPushLoginNotify = "StatSvc.SvcReqMSFLoginNotify"

	ServiceMethodConfigPushDomain   = "ConfigPushSvc.PushDomain"
	ServiceMethodConfigPushRequest  = "ConfigPushSvc.PushReq"
	ServiceMethodConfigPushResponse = "ConfigPushSvc.PushResp"

	ServiceMethodFriendListDeleteFriend       = "friendlist.delFriend"
	ServiceMethodFriendListGetFriendGroupList = "friendlist.getFriendGroupList"
	ServiceMethodFriendListGetGroupList       = "friendlist.GetTroopListReqV2"
	ServiceMethodFriendListGetGroupMemberList = "friendlist.GetTroopMemberListReq"

	ServiceMethodMessageDeleteMessage = "MessageSvc.PbDeleteMsg"
	ServiceMethodMessageGetMessage    = "MessageSvc.PbGetMsg"
	ServiceMethodMessageSendMessage   = "MessageSvc.PbSendMsg"
	ServiceMethodMessagePushNotify    = "MessageSvc.PushNotify"
	ServiceMethodMessagePushReaded    = "MessageSvc.PushReaded"

	ServiceMethodOnlinePushMessageSyncC2C   = "OnlinePush.PbC2CMsgSync"
	ServiceMethodOnlinePushMessageSyncGroup = "OnlinePush.PbPushGroupMsg"
	ServiceMethodOnlinePushMessageTransfer  = "OnlinePush.PbPushTransMsg"
	ServiceMethodOnlinePushRequest          = "OnlinePush.ReqPush"
	ServiceMethodOnlinePushResponse         = "OnlinePush.RespPush"
	ServiceMethodOnlinePushSIDTicketExpired = "OnlinePush.SidTicketExpired"
)
