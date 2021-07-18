package client

import (
	"context"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type FriendListGetGroupMemberListRequest struct {
	Uin                 int64 `jce:",0" json:",omitempty"`
	GroupCode           int64 `jce:",1" json:",omitempty"`
	NextUin             int64 `jce:",2" json:",omitempty"`
	GroupUin            int64 `jce:",3" json:",omitempty"`
	Version             int64 `jce:",4" json:",omitempty"`
	ReqType             int64 `jce:",5" json:",omitempty"`
	GetListAppointTime  int64 `jce:",6" json:",omitempty"`
	RichCardNameVersion int8  `jce:",7" json:",omitempty"`
}

type FriendListGetGroupMemberListResponse struct {
	Uin             int64             `jce:",0" json:",omitempty"`
	GroupCode       int64             `jce:",1" json:",omitempty"`
	GroupUin        int64             `jce:",2" json:",omitempty"`
	GroupMemberList []GroupMemberInfo `jce:",3" json:",omitempty"`
	NextUin         int64             `jce:",4" json:",omitempty"`
	Result          int32             `jce:",5" json:",omitempty"`
	ErrorCode       int16             `jce:",6" json:",omitempty"`
	OfficeMode      int64             `jce:",7" json:",omitempty"`
	NextGetTime     int64             `jce:",8" json:",omitempty"`
}

type GroupMemberInfo struct {
	MemberUin              int64          `jce:",0" json:",omitempty"`
	FaceID                 int16          `jce:",1" json:",omitempty"`
	Age                    int8           `jce:",2" json:",omitempty"`
	Gender                 int8           `jce:",3" json:",omitempty"`
	Nick                   string         `jce:",4" json:",omitempty"`
	Status                 int8           `jce:",5" json:",omitempty"`
	ShowName               string         `jce:",6" json:",omitempty"`
	Name                   string         `jce:",8" json:",omitempty"`
	Gender2                int8           `jce:",9" json:",omitempty"`
	Phone                  string         `jce:",10" json:",omitempty"`
	Email                  string         `jce:",11" json:",omitempty"`
	Memo                   string         `jce:",12" json:",omitempty"`
	AutoRemark             string         `jce:",13" json:",omitempty"`
	MemberLevel            int64          `jce:",14" json:",omitempty"`
	JoinTime               int64          `jce:",15" json:",omitempty"`
	LastSpeakTime          int64          `jce:",16" json:",omitempty"`
	CreditLevel            int64          `jce:",17" json:",omitempty"`
	Flag                   int64          `jce:",18" json:",omitempty"`
	FlagExt                int64          `jce:",19" json:",omitempty"`
	Point                  int64          `jce:",20" json:",omitempty"`
	Concerned              int8           `jce:",21" json:",omitempty"`
	Shielded               int8           `jce:",22" json:",omitempty"`
	SpecialTitle           string         `jce:",23" json:",omitempty"`
	SpecialTitleExpireTime int64          `jce:",24" json:",omitempty"`
	Job                    string         `jce:",25" json:",omitempty"`
	ApolloFlag             int8           `jce:",26" json:",omitempty"`
	ApolloTimestamp        int64          `jce:",27" json:",omitempty"`
	GlobalGroupLevel       int64          `jce:",28" json:",omitempty"`
	TitleId                int64          `jce:",29" json:",omitempty"`
	ShutupTimestap         int64          `jce:",30" json:",omitempty"`
	GlobalGroupPoint       int64          `jce:",31" json:",omitempty"`
	QZoneUserInfo          *QZoneUserInfo `jce:",32" json:",omitempty"`
	RichCardNameVer        int8           `jce:",33" json:",omitempty"`
	VipType                int64          `jce:",34" json:",omitempty"`
	VipLevel               int64          `jce:",35" json:",omitempty"`
	BigClubLevel           int64          `jce:",36" json:",omitempty"`
	BigClubFlag            int64          `jce:",37" json:",omitempty"`
	Nameplate              int64          `jce:",38" json:",omitempty"`
	GroupHonor             []byte         `jce:",39" json:",omitempty"`
	Remark                 []byte         `jce:",40" json:",omitempty"`
	RichFlag               int8           `jce:",41" json:",omitempty"`
}

type QZoneUserInfo struct {
	StarState  int32             `jce:",0" json:",omitempty"`
	ExtendInfo map[string]string `jce:",1" json:",omitempty"`
}

func NewFriendListGetGroupMemberListRequest(
	uin, groupCode, nextUin, groupUin int64,
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
		Username: strconv.FormatInt(req.Uin, 10),
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

	uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
	for _, item := range resp.GroupMemberList {
		channelMember := db.ChannelMember{
			Time:      item.JoinTime,
			Uin:       item.MemberUin,
			Nick:      item.Nick,
			Gender:    item.Gender + 1,
			Remark:    string(item.Remark),
			Status:    item.Status,
			ChannelID: resp.GroupCode,
		}
		if _, ok := c.cmembers[resp.GroupCode]; !ok {
			c.cmembers[resp.GroupCode] = make(map[int64]*db.ChannelMember)
		}
		c.cmembers[resp.GroupCode][item.MemberUin] = &channelMember
		if c.db != nil {
			err := c.dbInsertChannelMember(uin, &channelMember)
			if err != nil {
				log.Error().Err(err).Msg(">>> [db  ] dbInsertChannelMember")
			}
		}
	}
	return &resp, nil
}
