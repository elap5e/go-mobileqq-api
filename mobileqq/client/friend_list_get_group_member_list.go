package client

import (
	"context"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type FriendListGetGroupMemberListRequest struct {
	Uin                 uint64 `jce:",0"`
	GroupCode           uint64 `jce:",1"`
	NextUin             uint64 `jce:",2"`
	GroupUin            uint64 `jce:",3"`
	Version             uint64 `jce:",4"`
	ReqType             uint64 `jce:",5"`
	GetListAppointTime  uint64 `jce:",6"`
	RichCardNameVersion uint8  `jce:",7"`
}

type FriendListGetGroupMemberListResponse struct {
	Uin             uint64            `jce:",0"`
	GroupCode       uint64            `jce:",1"`
	GroupUin        uint64            `jce:",2"`
	GroupMemberList []GroupMemberInfo `jce:",3"`
	NextUin         uint64            `jce:",4"`
	Result          uint32            `jce:",5"`
	ErrorCode       uint16            `jce:",6"`
	OfficeMode      uint64            `jce:",7"`
	NextGetTime     uint64            `jce:",8"`
}

type GroupMemberInfo struct {
	MemberUin              uint64         `jce:",0"`
	FaceID                 uint16         `jce:",1"`
	Age                    uint8          `jce:",2"`
	Gender                 uint8          `jce:",3"`
	Nick                   string         `jce:",4"`
	Status                 uint8          `jce:",5"`
	ShowName               string         `jce:",6"`
	Name                   string         `jce:",8"`
	Gender2                uint8          `jce:",9"`
	Phone                  string         `jce:",10"`
	Email                  string         `jce:",11"`
	Memo                   string         `jce:",12"`
	AutoRemark             string         `jce:",13"`
	MemberLevel            uint64         `jce:",14"`
	JoinTime               uint64         `jce:",15"`
	LastSpeakTime          uint64         `jce:",16"`
	CreditLevel            uint64         `jce:",17"`
	Flag                   uint64         `jce:",18"`
	FlagExt                uint64         `jce:",19"`
	Point                  uint64         `jce:",20"`
	Concerned              uint8          `jce:",21"`
	Shielded               uint8          `jce:",22"`
	SpecialTitle           string         `jce:",23"`
	SpecialTitleExpireTime uint64         `jce:",24"`
	Job                    string         `jce:",25"`
	ApolloFlag             uint8          `jce:",26"`
	ApolloTimestamp        uint64         `jce:",27"`
	GlobalGroupLevel       uint64         `jce:",28"`
	TitleId                uint64         `jce:",29"`
	ShutupTimestap         uint64         `jce:",30"`
	GlobalGroupPoint       uint64         `jce:",31"`
	QZoneUserInfo          *QZoneUserInfo `jce:",32"`
	RichCardNameVer        uint8          `jce:",33"`
	VipType                uint64         `jce:",34"`
	VipLevel               uint64         `jce:",35"`
	BigClubLevel           uint64         `jce:",36"`
	BigClubFlag            uint64         `jce:",37"`
	Nameplate              uint64         `jce:",38"`
	GroupHonor             []byte         `jce:",39"`
	Name2                  []byte         `jce:",40"`
	RichFlag               uint8          `jce:",41"`
}

type QZoneUserInfo struct {
	StarState  uint32            `jce:",0"`
	ExtendInfo map[string]string `jce:",1"`
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
