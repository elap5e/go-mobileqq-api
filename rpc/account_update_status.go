package rpc

import (
	"context"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
)

type PushRegisterInfoStatusType uint32

var (
	PushRegisterInfoStatusOnline                PushRegisterInfoStatusType = 0x0000000b // 11
	PushRegisterInfoStatusOffline               PushRegisterInfoStatusType = 0x00000015 // 21
	PushRegisterInfoStatusAway                  PushRegisterInfoStatusType = 0x0000001f // 31
	PushRegisterInfoStatusInvisiable            PushRegisterInfoStatusType = 0x00000029 // 41
	PushRegisterInfoStatusReceiveOfflineMessage PushRegisterInfoStatusType = 0x0000005f // 95
)

type AccountUpdateStatus struct {
	Uin       uint64   `jce:",1"`
	PushIDs   []uint64 `jce:",2"` // constant
	Status    uint32   `jce:",3"`
	KikPC     bool     `jce:",4"`
	KikWeak   bool     `jce:",5"` // constant false
	Timestamp uint64   `jce:",6"`
	LargeSeq  uint32   `jce:",7"` // constant 0x00000000
}

type AccountUpdateStatusRequest struct {
	Uin          uint64 `jce:",0"`
	Bid          uint64 `jce:",1"`
	ConnType     uint8  `jce:",2"` // constant 0x00
	Other        string `jce:",3"` // constant ""
	Status       uint32 `jce:",4"`
	OnlinePush   bool   `jce:",5"` // constant false
	IsOnline     bool   `jce:",6"` // constant false
	IsShowOnline bool   `jce:",7"` // constant false
	KikPC        bool   `jce:",8"`
	KikWeak      bool   `jce:",9"` // constant false
	Timestamp    uint64 `jce:",10"`
	SDKVersion   uint32 `jce:",11"`
	NetworkType  uint8  `jce:",12"` // 0x00: mobile; 0x01: wifi
	BuildVersion string `jce:",13"` // constant ""
	RegisterType bool   `jce:",14"` // false: appRegister, fillRegProxy, createDefaultRegInfo; true: others
	DevParam     []byte `jce:",15"` // constant nil
	GUID         []byte `jce:",16"` // placeholder
	LocaleID     uint32 `jce:",17"` // constant 0x00000804
	SlientPush   bool   `jce:",18"` // constant false
	DeviceName   string `jce:",19"`
	DeviceType   string `jce:",20"`
	OSVersion    string `jce:",21"`
	OpenPush     bool   `jce:",22"` // constant true
	LargeSeq     uint32 `jce:",23"` // constant 0x00000000

	// LastWatchStartTime uint32         `jce:",24"`
	// BindUin            []uint64       `jce:",25"`
	// OldSSOIP           uint64         `jce:",26"`
	// NewSSOIP           uint64         `jce:",27"`
	// ChannelNo          string         `jce:",28"`
	// CPID               uint64         `jce:",29"`
	// VendorName         string         `jce:",30"`
	// VendorOSName       string         `jce:",31"`
	// IOSIDFA            string         `jce:",32"`
	// Reqbody0x769       []byte         `jce:",33"`
	// IsSetStatus        bool           `jce:",34"`
	// ServerBuf          []byte         `jce:",35"`
	// SetMute            bool           `jce:",36"`
	// NotifySwitch       uint8          `jce:",37"`
	// ExtOnlineStatus    uint64         `jce:",38"`
	// BatteryStatus      uint32         `jce:",39"`
	// VendorPushInfo     VendorPushInfo `jce:",42"`
}

type VendorPushInfo struct {
	Type uint64 `jce:",0"`
}

type AccountUpdateStatusResponse struct {
	Uin            uint64 `jce:",0"`
	Bid            uint64 `jce:",1"`
	ReplyCode      uint8  `jce:",2"`
	Result         string `jce:",3"`
	ServerTime     uint64 `jce:",4"`
	LogQQ          bool   `jce:",5"`
	NeedKik        bool   `jce:",6"`
	UpdateFlag     bool   `jce:",7"`
	Timestamp      uint64 `jce:",8"`
	CrashFlag      bool   `jce:",9"`
	ClientIP       string `jce:",10"`
	ClientPort     uint32 `jce:",11"`
	HelloInterval  uint32 `jce:",12"`
	LargeSeq       uint32 `jce:",13"`
	LargeSeqUpdate bool   `jce:",14"`

	Respbody0x769            []byte `jce:",15"`
	Status                   uint32 `jce:",16"`
	ExtraOnlineStatus        uint64 `jce:",17"`
	ClientBatteryGetInterval uint64 `jce:",18"`
	ClientAutoStatusInterval uint64 `jce:",19"`
}

func NewAccountUpdateStatusRequest(
	uin uint64,
	status PushRegisterInfoStatusType,
	kick bool,
) *AccountUpdateStatusRequest {
	ids := []uint64{0x01, 0x02, 0x04}
	bid := uint64(0x0000000000000000)
	for _, id := range ids {
		bid |= id
	}
	push := &AppPushInfo{
		Bid: bid,
		AccountUpdateStatus: AccountUpdateStatus{
			Uin:       uin,
			PushIDs:   ids,
			Status:    uint32(status),
			KikPC:     kick,
			KikWeak:   false,
			Timestamp: 0x0000000000000000, // TODO: fix
			LargeSeq:  0x00000000,
		},
	}
	return &AccountUpdateStatusRequest{
		Uin:          push.AccountUpdateStatus.Uin,
		Bid:          push.Bid,
		ConnType:     0x00,
		Other:        "",
		Status:       push.AccountUpdateStatus.Status,
		OnlinePush:   false,
		IsOnline:     false,
		IsShowOnline: false,
		KikPC:        push.AccountUpdateStatus.KikPC,
		KikWeak:      push.AccountUpdateStatus.KikWeak,
		Timestamp:    push.AccountUpdateStatus.Timestamp,
		SDKVersion:   defaultDeviceOSSDKVersion,
		NetworkType:  0x01,
		BuildVersion: "",
		RegisterType: false,
		DevParam:     nil,
		GUID:         nil,
		LocaleID:     0x00000804,
		SlientPush:   false,
		DeviceName:   defaultDeviceOSBuildModel,
		DeviceType:   defaultDeviceOSBuildModel,
		OSVersion:    defaultDeviceOSVersion,
		OpenPush:     true,
		LargeSeq:     push.AccountUpdateStatus.LargeSeq,
	}
}

func (c *Client) AccountUpdateStatus(
	ctx context.Context,
	req *AccountUpdateStatusRequest,
) (*AccountUpdateStatusResponse, error) {
	req.GUID = c.cfg.Device.GUID
	buf, err := uni.Marshal(ctx, &uni.Message{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   0x00000000,
		ServantName: "PushService",
		FuncName:    "SvcReqRegister",
		Buffer:      []byte{},
		Timeout:     0x00000000,
		Context:     map[string]string{},
		Status:      map[string]string{},
	}, map[string]interface{}{
		"SvcReqRegister": req,
	})
	if err != nil {
		return nil, err
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call(ServiceMethodAccountUpdateStatus, &ClientToServerMessage{
		Username: strconv.FormatInt(int64(req.Uin), 10),
		Seq:      c.getNextSeq(),
		Buffer:   buf,
		Simple:   false,
	}, s2c); err != nil {
		return nil, err
	}
	msg := new(uni.Message)
	resp := new(AccountUpdateStatusResponse)
	if err := uni.Unmarshal(ctx, s2c.Buffer, msg, map[string]interface{}{
		"SvcRespRegister": resp,
	}); err != nil {
		return nil, err
	}
	return resp, nil
}
