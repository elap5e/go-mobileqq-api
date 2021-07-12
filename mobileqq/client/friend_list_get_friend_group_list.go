package client

import (
	"context"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type FriendListGetFriendGroupListRequest struct {
	RequestType     uint32   `jce:",0" json:",omitempty"`
	IsReflush       bool     `jce:",1" json:",omitempty"`
	Uin             uint64   `jce:",2" json:",omitempty"`
	StartIndex      uint16   `jce:",3" json:",omitempty"`
	GetFriendCount  uint16   `jce:",4" json:",omitempty"`
	GroupID         uint8    `jce:",5" json:",omitempty"`
	IsGetGroupInfo  bool     `jce:",6" json:",omitempty"`
	GroupStartIndex uint8    `jce:",7" json:",omitempty"`
	GetGroupCount   uint8    `jce:",8" json:",omitempty"`
	IsGetMSFGroup   bool     `jce:",9" json:",omitempty"`
	IsShowTermType  bool     `jce:",10" json:",omitempty"`
	Version         uint64   `jce:",11" json:",omitempty"`
	UinList         []uint64 `jce:",12" json:",omitempty"`
	AppType         uint32   `jce:",13" json:",omitempty"`
	IsGetDOVID      bool     `jce:",14" json:",omitempty"`
	IsGetBothFlag   bool     `jce:",15" json:",omitempty"`
	OIDB0x0D50      []byte   `jce:",16" json:",omitempty"`
	OIDB0x0D6B      []byte   `jce:",17" json:",omitempty"`
	SNSTypeList     []uint64 `jce:",18" json:",omitempty"`
}

type FriendListGetFriendGroupListResponse struct {
	RequestType             uint32                           `jce:",0" json:",omitempty"`
	IsReflush               bool                             `jce:",1" json:",omitempty"`
	Uin                     uint64                           `jce:",2" json:",omitempty"`
	StartIndex              uint16                           `jce:",3" json:",omitempty"`
	GetFriendCount          uint16                           `jce:",4" json:",omitempty"`
	TotoalFriendCount       uint16                           `jce:",5" json:",omitempty"`
	FriendCount             uint16                           `jce:",6" json:",omitempty"`
	FriendInfoList          []FriendInfo                     `jce:",7" json:",omitempty"`
	GroupID                 uint8                            `jce:",8" json:",omitempty"`
	IsGetGroupInfo          bool                             `jce:",9" json:",omitempty"`
	GroupStartIndex         uint8                            `jce:",10" json:",omitempty"`
	GetGroupCount           uint8                            `jce:",11" json:",omitempty"`
	TotoalGroupCount        uint16                           `jce:",12" json:",omitempty"`
	GroupCount              uint8                            `jce:",13" json:",omitempty"`
	GroupInfoList           []FriendGroupInfo                `jce:",14" json:",omitempty"`
	Result                  uint32                           `jce:",15" json:",omitempty"`
	ErrorCode               uint16                           `jce:",16" json:",omitempty"`
	OnlineFriendCount       uint16                           `jce:",17" json:",omitempty"`
	ServerTime              uint64                           `jce:",18" json:",omitempty"`
	QQOnlineCount           uint16                           `jce:",19" json:",omitempty"`
	GroupInfoList2          []FriendGroupInfo                `jce:",20" json:",omitempty"`
	RespType                uint8                            `jce:",21" json:",omitempty"`
	HasOtherRespFlag        uint8                            `jce:",22" json:",omitempty"`
	FriendInfo              *FriendInfo                      `jce:",23" json:",omitempty"`
	ShowPcIcon              uint8                            `jce:",24" json:",omitempty"`
	GetExtraSNSResponseCode uint16                           `jce:",25" json:",omitempty"`
	SubServerResponseCode   *FriendListSubServerResponseCode `jce:",26" json:",omitempty"`
}

type FriendInfo struct {
	FriendUin             uint64       `jce:",0" json:",omitempty"`
	GroupID               uint8        `jce:",1" json:",omitempty"`
	FaceID                uint16       `jce:",2" json:",omitempty"`
	Remark                string       `jce:",3" json:",omitempty"`
	QQType                uint8        `jce:",4" json:",omitempty"`
	Status                uint8        `jce:",5" json:",omitempty"`
	MemberLevel           uint8        `jce:",6" json:",omitempty"`
	IsMobileQQOnLine      bool         `jce:",7" json:",omitempty"`
	QQOnLineState         uint8        `jce:",8" json:",omitempty"`
	IsIphoneOnline        bool         `jce:",9" json:",omitempty"`
	DetalStatusFlag       uint8        `jce:",10" json:",omitempty"`
	QQOnLineStateV2       uint8        `jce:",11" json:",omitempty"`
	ShowName              string       `jce:",12" json:",omitempty"`
	IsRemark              bool         `jce:",13" json:",omitempty"`
	Nick                  string       `jce:",14" json:",omitempty"`
	SpecialFlag           uint8        `jce:",15" json:",omitempty"`
	IMGroupID             []byte       `jce:",16" json:",omitempty"`
	MSFGroupID            []byte       `jce:",17" json:",omitempty"`
	TermType              uint32       `jce:",18" json:",omitempty"`
	VIPBaseInfo           *VIPBaseInfo `jce:",19" json:",omitempty"`
	Network               uint8        `jce:",20" json:",omitempty"`
	Ring                  []byte       `jce:",21" json:",omitempty"`
	AbiFlag               uint64       `jce:",22" json:",omitempty"`
	FaceAddonId           uint64       `jce:",23" json:",omitempty"`
	NetworkType           uint32       `jce:",24" json:",omitempty"`
	VIPFont               uint64       `jce:",25" json:",omitempty"`
	IconType              uint32       `jce:",26" json:",omitempty"`
	TermDesc              string       `jce:",27" json:",omitempty"`
	ColorRing             uint64       `jce:",28" json:",omitempty"`
	ApolloFlag            uint8        `jce:",29" json:",omitempty"`
	ApolloTimestamp       uint64       `jce:",30" json:",omitempty"`
	Gender                uint8        `jce:",31" json:",omitempty"`
	FounderFont           uint64       `jce:",32" json:",omitempty"`
	EimId                 string       `jce:",33" json:",omitempty"`
	EimMobile             string       `jce:",34" json:",omitempty"`
	OlympicTorch          uint8        `jce:",35" json:",omitempty"`
	ApolloSignTime        uint64       `jce:",36" json:",omitempty"`
	LaviUin               uint64       `jce:",37" json:",omitempty"`
	TagUpdateTime         uint64       `jce:",38" json:",omitempty"`
	GameLastLoginTime     uint64       `jce:",39" json:",omitempty"`
	GameAppID             uint64       `jce:",40" json:",omitempty"`
	CardID                []byte       `jce:",41" json:",omitempty"`
	BitSet                uint64       `jce:",42" json:",omitempty"`
	KingOfGloryFlag       uint8        `jce:",43" json:",omitempty"`
	KingOfGloryRank       uint64       `jce:",44" json:",omitempty"`
	MasterUin             string       `jce:",45" json:",omitempty"`
	LastMedalUpdateTime   uint64       `jce:",46" json:",omitempty"`
	FaceStoreId           uint64       `jce:",47" json:",omitempty"`
	FontEffect            uint64       `jce:",48" json:",omitempty"`
	DOVID                 string       `jce:",49" json:",omitempty"`
	BothFlag              uint64       `jce:",50" json:",omitempty"`
	CentiShow3DFlag       uint8        `jce:",51" json:",omitempty"`
	IntimateInfo          []byte       `jce:",52" json:",omitempty"`
	ShowNameplate         uint8        `jce:",53" json:",omitempty"`
	NewLoverDiamondFlag   uint8        `jce:",54" json:",omitempty"`
	ExtSnsFrdData         []byte       `jce:",55" json:",omitempty"`
	MutualMarkData        []byte       `jce:",56" json:",omitempty"`
	ExtOnlineStatus       uint64       `jce:",57" json:",omitempty"`
	BatteryStatus         uint32       `jce:",58" json:",omitempty"`
	MusicInfo             []byte       `jce:",59" json:",omitempty"`
	PoiInfo               []byte       `jce:",60" json:",omitempty"`
	ExtOnlineBusinessInfo []byte       `jce:",61" json:",omitempty"`
}

type VIPBaseInfo struct {
	OpenInfoMap       map[uint64]VIPOpenInfo `jce:",0" json:",omitempty"`
	NameplateVIPType  uint32                 `jce:",1" json:",omitempty"`
	GrayNameplateFlag uint32                 `jce:",2" json:",omitempty"`
	ExtendNameplateId string                 `jce:",3" json:",omitempty"`
}

type VIPOpenInfo struct {
	Open        bool   `jce:",0" json:",omitempty"`
	VIPType     uint32 `jce:",1" json:",omitempty"`
	VIPLevel    uint32 `jce:",2" json:",omitempty"`
	VIPFlag     uint32 `jce:",3" json:",omitempty"`
	NameplateID uint64 `jce:",4" json:",omitempty"`
}

type FriendListSubServerResponseCode struct {
	GetMutualMarkRspCode   uint16 `jce:",0" json:",omitempty"`
	GetIntimateInfoRspCode uint16 `jce:",1" json:",omitempty"`
}

type FriendGroupInfo struct {
	GroupID   uint8  `jce:",0" json:",omitempty"`
	GroupName string `jce:",1" json:",omitempty"`
}

func NewFriendListGetFriendGroupListRequest(
	uin uint64,
	startIndex, friendCount uint16,
	groupStartIndex, groupCount uint8,
) *FriendListGetFriendGroupListRequest {
	oidb0x0D50, _ := proto.Marshal(&pb.OIDB0X0D50Request{
		AppId:                   0x000000000002712,
		ReqMusicSwitch:          0x00000001,
		ReqKsingSwitch:          0x00000001,
		ReqMutualmarkLbsshare:   0x00000001,
		ReqMutualmarkAlienation: 0x00000001,
		ReqAioQuickApp:          0x00000001,
	})
	oidb0x0D6B, _ := proto.Marshal(&pb.OIDB0X0D6BRequest{})
	return &FriendListGetFriendGroupListRequest{
		RequestType:     0x00000003,
		IsReflush:       startIndex == 0,
		Uin:             uin,
		StartIndex:      startIndex,
		GetFriendCount:  friendCount,
		GroupID:         0x00,
		IsGetGroupInfo:  startIndex == 0,
		GroupStartIndex: groupStartIndex,
		GetGroupCount:   groupCount,
		IsGetMSFGroup:   false,
		IsShowTermType:  true,
		Version:         0x000000000000001f,
		UinList:         nil,
		AppType:         0x00000000,
		IsGetDOVID:      false,
		IsGetBothFlag:   false,
		OIDB0x0D50:      oidb0x0D50,
		OIDB0x0D6B:      oidb0x0D6B,
		SNSTypeList: []uint64{
			0x000000000000350c,
			0x000000000000350d,
			0x000000000000350e,
		},
	}
}

func (c *Client) FriendListGetFriendGroupList(
	ctx context.Context,
	req *FriendListGetFriendGroupListRequest,
) (*FriendListGetFriendGroupListResponse, error) {
	buf, err := uni.Marshal(ctx, &uni.Message{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   c.getNextRequestSeq(),
		ServantName: "mqq.IMService.FriendListServiceServantObj",
		FuncName:    "GetFriendListReq",
		Buffer:      []byte{},
		Timeout:     0x00000000,
		Context:     map[string]string{},
		Status:      map[string]string{},
	}, map[string]interface{}{
		"FL": req,
	})
	if err != nil {
		return nil, err
	}
	c2s, s2c := codec.ClientToServerMessage{
		Username: strconv.FormatInt(int64(req.Uin), 10),
		Buffer:   buf,
		Simple:   true,
	}, codec.ServerToClientMessage{}
	err = c.rpc.Call(ServiceMethodFriendListGetFriendGroupList, &c2s, &s2c)
	if err != nil {
		return nil, err
	}
	msg := uni.Message{}
	resp := FriendListGetFriendGroupListResponse{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"FLRESP": &resp,
	}); err != nil {
		return nil, err
	}

	dumpServerToClientMessage(&s2c, &resp)

	for i := range resp.FriendInfoList {
		c.contacts[resp.FriendInfoList[i].FriendUin] = &resp.FriendInfoList[i]
	}
	return &resp, nil
}
