package client

import (
	"context"
	"strconv"
	"time"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type AccountStatusType uint32

var (
	AccountStatusOnline                AccountStatusType = 0x0000000b // 11
	AccountStatusOffline               AccountStatusType = 0x00000015 // 21
	AccountStatusAway                  AccountStatusType = 0x0000001f // 31
	AccountStatusInvisiable            AccountStatusType = 0x00000029 // 41
	AccountStatusReceiveOfflineMessage AccountStatusType = 0x0000005f // 95
)

type AccountStatus struct {
	Uin       uint64   `jce:",1" json:",omitempty"`
	PushIDs   []uint64 `jce:",2" json:",omitempty"` // constant
	Status    uint32   `jce:",3" json:",omitempty"`
	KickPC    bool     `jce:",4" json:",omitempty"`
	KickWeak  bool     `jce:",5" json:",omitempty"` // constant false
	Timestamp uint64   `jce:",6" json:",omitempty"`
	LargeSeq  uint32   `jce:",7" json:",omitempty"` // constant 0x00000000
}

type AccountSetStatusRequest struct {
	Uin          uint64 `jce:",0" json:",omitempty"`
	Bid          uint64 `jce:",1" json:",omitempty"`
	ConnType     uint8  `jce:",2" json:",omitempty"` // constant 0x00
	Other        string `jce:",3" json:",omitempty"` // constant ""
	Status       uint32 `jce:",4" json:",omitempty"`
	OnlinePush   bool   `jce:",5" json:",omitempty"` // constant false
	IsOnline     bool   `jce:",6" json:",omitempty"` // constant false
	IsShowOnline bool   `jce:",7" json:",omitempty"` // constant false
	KickPC       bool   `jce:",8" json:",omitempty"`
	KickWeak     bool   `jce:",9" json:",omitempty"` // constant false
	Timestamp    uint64 `jce:",10" json:",omitempty"`
	SDKVersion   uint32 `jce:",11" json:",omitempty"`
	NetworkType  uint8  `jce:",12" json:",omitempty"` // 0x00: mobile; 0x01: wifi
	BuildVersion string `jce:",13" json:",omitempty"` // constant ""
	RegisterType bool   `jce:",14" json:",omitempty"` // false: appRegister, fillRegProxy, createDefaultRegInfo; true: others
	DevParam     []byte `jce:",15" json:",omitempty"` // constant nil
	GUID         []byte `jce:",16" json:",omitempty"` // placeholder
	LocaleID     uint32 `jce:",17" json:",omitempty"` // constant 0x00000804
	SlientPush   bool   `jce:",18" json:",omitempty"` // constant false
	DeviceName   string `jce:",19" json:",omitempty"`
	DeviceType   string `jce:",20" json:",omitempty"`
	OSVersion    string `jce:",21" json:",omitempty"`
	OpenPush     bool   `jce:",22" json:",omitempty"` // constant true
	LargeSeq     uint32 `jce:",23" json:",omitempty"` // constant 0x00000000

	// LastWatchStartTime uint32         `jce:",24" json:",omitempty"`
	// BindUin            []uint64       `jce:",25" json:",omitempty"`
	// OldSSOIP           uint64         `jce:",26" json:",omitempty"`
	// NewSSOIP           uint64         `jce:",27" json:",omitempty"`
	// ChannelNo          string         `jce:",28" json:",omitempty"`
	// CPID               uint64         `jce:",29" json:",omitempty"`
	// VendorName         string         `jce:",30" json:",omitempty"`
	// VendorOSName       string         `jce:",31" json:",omitempty"`
	// IOSIDFA            string         `jce:",32" json:",omitempty"`
	// Reqbody0x769       []byte         `jce:",33" json:",omitempty"`
	// IsSetStatus        bool           `jce:",34" json:",omitempty"`
	// ServerBuf          []byte         `jce:",35" json:",omitempty"`
	// SetMute            bool           `jce:",36" json:",omitempty"`
	// NotifySwitch       uint8          `jce:",37" json:",omitempty"`
	// ExtOnlineStatus    uint64         `jce:",38" json:",omitempty"`
	// BatteryStatus      uint32         `jce:",39" json:",omitempty"`
	// VendorPushInfo     VendorPushInfo `jce:",42" json:",omitempty"`
}

type VendorPushInfo struct {
	Type uint64 `jce:",0" json:",omitempty"`
}

type AccountSetStatusResponse struct {
	Uin            uint64 `jce:",0" json:",omitempty"`
	Bid            uint64 `jce:",1" json:",omitempty"`
	ReplyCode      uint8  `jce:",2" json:",omitempty"`
	Result         string `jce:",3" json:",omitempty"`
	ServerTime     uint64 `jce:",4" json:",omitempty"`
	LogQQ          bool   `jce:",5" json:",omitempty"`
	NeedKick       bool   `jce:",6" json:",omitempty"`
	UpdateFlag     bool   `jce:",7" json:",omitempty"`
	Timestamp      uint64 `jce:",8" json:",omitempty"`
	CrashFlag      bool   `jce:",9" json:",omitempty"`
	ClientIP       string `jce:",10" json:",omitempty"`
	ClientPort     uint32 `jce:",11" json:",omitempty"`
	HelloInterval  uint32 `jce:",12" json:",omitempty"`
	LargeSeq       uint32 `jce:",13" json:",omitempty"`
	LargeSeqUpdate bool   `jce:",14" json:",omitempty"`

	Respbody0x769            []byte `jce:",15" json:",omitempty"`
	Status                   uint32 `jce:",16" json:",omitempty"`
	ExtraOnlineStatus        uint64 `jce:",17" json:",omitempty"`
	ClientBatteryGetInterval uint64 `jce:",18" json:",omitempty"`
	ClientAutoStatusInterval uint64 `jce:",19" json:",omitempty"`
}

func (c *Client) initAutoStatusTimers() {
	c.autoStatusTimers = make(map[string]*time.Timer)
}

func NewAccountSetStatusRequest(
	uin uint64,
	status AccountStatusType,
	kick bool,
) *AccountSetStatusRequest {
	ids := []uint64{0x01, 0x02, 0x04}
	bid := uint64(0x0000000000000000)
	for _, id := range ids {
		bid |= id
	}
	push := &AppPushInfo{
		Bid: bid,
		AccountStatus: AccountStatus{
			Uin:       uin,
			PushIDs:   ids,
			Status:    uint32(status),
			KickPC:    kick,
			KickWeak:  false,
			Timestamp: 0x0000000000000000, // TODO: fix
			LargeSeq:  0x00000000,
		},
	}
	return &AccountSetStatusRequest{
		Uin:          push.AccountStatus.Uin,
		Bid:          push.Bid,
		ConnType:     0x00,
		Other:        "",
		Status:       push.AccountStatus.Status,
		OnlinePush:   false,
		IsOnline:     false,
		IsShowOnline: false,
		KickPC:       push.AccountStatus.KickPC,
		KickWeak:     push.AccountStatus.KickWeak,
		Timestamp:    push.AccountStatus.Timestamp,
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
		LargeSeq:     push.AccountStatus.LargeSeq,
	}
}

func (c *Client) AccountSetStatus(
	ctx context.Context,
	req *AccountSetStatusRequest,
) (*AccountSetStatusResponse, error) {
	req.GUID = c.cfg.Device.GUID
	buf, err := uni.Marshal(ctx, &uni.Message{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   c.getNextRequestSeq(),
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
	c2s, s2c := codec.ClientToServerMessage{
		Username: strconv.FormatInt(int64(req.Uin), 10),
		Buffer:   buf,
		Simple:   false,
	}, codec.ServerToClientMessage{}
	err = c.rpc.Call(ServiceMethodAccountSetStatus, &c2s, &s2c)
	if err != nil {
		return nil, err
	}
	msg := uni.Message{}
	resp := AccountSetStatusResponse{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"SvcRespRegister": &resp,
	}); err != nil {
		return nil, err
	}

	dumpServerToClientMessage(&s2c, &resp)

	c.autoStatusTimersMux.Lock()
	timer, ok := c.autoStatusTimers[c2s.Username]
	if !ok {
		c.autoStatusTimers[c2s.Username] = time.AfterFunc(0, func() {
			c.AccountSetStatus(ctx, req)
		})
		timer = c.autoStatusTimers[c2s.Username]
	}
	if resp.ClientAutoStatusInterval > 30 {
		timer.Reset(time.Duration(resp.ClientAutoStatusInterval) * time.Second)
	} else {
		timer.Reset(30 * time.Second)
	}
	c.autoStatusTimersMux.Unlock()
	return &resp, nil
}
