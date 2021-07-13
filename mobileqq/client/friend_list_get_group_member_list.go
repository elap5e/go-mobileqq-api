package client

import (
	"context"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type FriendListGetGroupMemberListRequest struct {
	Uin                 uint64 `jce:",0" json:",omitempty"`
	GroupCode           uint64 `jce:",1" json:",omitempty"`
	NextUin             uint64 `jce:",2" json:",omitempty"`
	GroupUin            uint64 `jce:",3" json:",omitempty"`
	Version             uint64 `jce:",4" json:",omitempty"`
	ReqType             uint64 `jce:",5" json:",omitempty"`
	GetListAppointTime  uint64 `jce:",6" json:",omitempty"`
	RichCardNameVersion uint8  `jce:",7" json:",omitempty"`
}

type FriendListGetGroupMemberListResponse struct {
	Uin             uint64            `jce:",0" json:",omitempty"`
	GroupCode       uint64            `jce:",1" json:",omitempty"`
	GroupUin        uint64            `jce:",2" json:",omitempty"`
	GroupMemberList []GroupMemberInfo `jce:",3" json:",omitempty"`
	NextUin         uint64            `jce:",4" json:",omitempty"`
	Result          uint32            `jce:",5" json:",omitempty"`
	ErrorCode       uint16            `jce:",6" json:",omitempty"`
	OfficeMode      uint64            `jce:",7" json:",omitempty"`
	NextGetTime     uint64            `jce:",8" json:",omitempty"`
}

type GroupMemberInfo struct {
	MemberUin              uint64         `jce:",0" json:",omitempty"`
	FaceID                 uint16         `jce:",1" json:",omitempty"`
	Age                    uint8          `jce:",2" json:",omitempty"`
	Gender                 uint8          `jce:",3" json:",omitempty"`
	Nick                   string         `jce:",4" json:",omitempty"`
	Status                 uint8          `jce:",5" json:",omitempty"`
	ShowName               string         `jce:",6" json:",omitempty"`
	Name                   string         `jce:",8" json:",omitempty"`
	Gender2                uint8          `jce:",9" json:",omitempty"`
	Phone                  string         `jce:",10" json:",omitempty"`
	Email                  string         `jce:",11" json:",omitempty"`
	Memo                   string         `jce:",12" json:",omitempty"`
	AutoRemark             string         `jce:",13" json:",omitempty"`
	MemberLevel            uint64         `jce:",14" json:",omitempty"`
	JoinTime               uint64         `jce:",15" json:",omitempty"`
	LastSpeakTime          uint64         `jce:",16" json:",omitempty"`
	CreditLevel            uint64         `jce:",17" json:",omitempty"`
	Flag                   uint64         `jce:",18" json:",omitempty"`
	FlagExt                uint64         `jce:",19" json:",omitempty"`
	Point                  uint64         `jce:",20" json:",omitempty"`
	Concerned              uint8          `jce:",21" json:",omitempty"`
	Shielded               uint8          `jce:",22" json:",omitempty"`
	SpecialTitle           string         `jce:",23" json:",omitempty"`
	SpecialTitleExpireTime uint64         `jce:",24" json:",omitempty"`
	Job                    string         `jce:",25" json:",omitempty"`
	ApolloFlag             uint8          `jce:",26" json:",omitempty"`
	ApolloTimestamp        uint64         `jce:",27" json:",omitempty"`
	GlobalGroupLevel       uint64         `jce:",28" json:",omitempty"`
	TitleId                uint64         `jce:",29" json:",omitempty"`
	ShutupTimestap         uint64         `jce:",30" json:",omitempty"`
	GlobalGroupPoint       uint64         `jce:",31" json:",omitempty"`
	QZoneUserInfo          *QZoneUserInfo `jce:",32" json:",omitempty"`
	RichCardNameVer        uint8          `jce:",33" json:",omitempty"`
	VipType                uint64         `jce:",34" json:",omitempty"`
	VipLevel               uint64         `jce:",35" json:",omitempty"`
	BigClubLevel           uint64         `jce:",36" json:",omitempty"`
	BigClubFlag            uint64         `jce:",37" json:",omitempty"`
	Nameplate              uint64         `jce:",38" json:",omitempty"`
	GroupHonor             []byte         `jce:",39" json:",omitempty"`
	Name2                  []byte         `jce:",40" json:",omitempty"`
	RichFlag               uint8          `jce:",41" json:",omitempty"`
}

type QZoneUserInfo struct {
	StarState  uint32            `jce:",0" json:",omitempty"`
	ExtendInfo map[string]string `jce:",1" json:",omitempty"`
}

func NewFriendListGetGroupMemberListRequest(
	uin, groupCode, nextUin, groupUin uint64,
) *FriendListGetGroupMemberListRequest {
	return &FriendListGetGroupMemberListRequest{
		Uin:                 uin,
		GroupCode:           groupCode,
		NextUin:             nextUin,
		GroupUin:            groupUin,
		Version:             0x0000000000000003,
		ReqType:             0x0000000000000000,
		GetListAppointTime:  0x0000000000000000,
		RichCardNameVersion: 0x01,
	}
}

func (c *Client) FriendListGetGroupMemberList(
	ctx context.Context,
	req *FriendListGetGroupMemberListRequest,
) (*FriendListGetGroupMemberListResponse, error) {
	buf, err := uni.Marshal(ctx, &uni.Message{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   c.getNextRequestSeq(),
		ServantName: "mqq.IMService.FriendListServiceServantObj",
		FuncName:    "GetTroopMemberListReq",
		Buffer:      []byte{},
		Timeout:     0x00000000,
		Context:     map[string]string{},
		Status:      map[string]string{},
	}, map[string]interface{}{
		"GTML": req,
	})
	if err != nil {
		return nil, err
	}
	c2s, s2c := codec.ClientToServerMessage{
		Username: strconv.FormatInt(int64(req.Uin), 10),
		Buffer:   buf,
		Simple:   true,
	}, codec.ServerToClientMessage{}
	err = c.rpc.Call(ServiceMethodFriendListGetGroupMemberList, &c2s, &s2c)
	if err != nil {
		return nil, err
	}
	msg := uni.Message{}
	resp := FriendListGetGroupMemberListResponse{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"GTMLRESP": &resp,
	}); err != nil {
		return nil, err
	}

	dumpServerToClientMessage(&s2c, &resp)

	for i := range resp.GroupMemberList {
		if _, ok := c.cmembers[resp.GroupCode]; !ok {
			c.cmembers[resp.GroupCode] = make(map[uint64]*GroupMemberInfo)
		}
		c.cmembers[resp.GroupCode][resp.GroupMemberList[i].MemberUin] = &resp.GroupMemberList[i]
	}
	return &resp, nil
}
