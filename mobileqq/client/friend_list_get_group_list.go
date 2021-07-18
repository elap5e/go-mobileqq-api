package client

import (
	"context"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type FriendListGetGroupListRequest struct {
	Uin               int64              `jce:",0" json:",omitempty"`
	GetMSFMessageFlag int8               `jce:",1" json:",omitempty"`
	Cookie            []byte             `jce:",2" json:",omitempty"`
	GroupInfoList     []GroupInfoRequest `jce:",3" json:",omitempty"`
	GroupFlagExtra    int8               `jce:",4" json:",omitempty"`
	Version           int32              `jce:",5" json:",omitempty"`
	CompanyID         int64              `jce:",6" json:",omitempty"`
	VersionNumber     int64              `jce:",7" json:",omitempty"`
	GetLongGroupName  int8               `jce:",8" json:",omitempty"`
}

type GroupInfoRequest struct {
	GroupCode         int64 `jce:",0" json:",omitempty"`
	GroupInfoSeq      int64 `jce:",1" json:",omitempty"`
	GroupFlagExtra    int64 `jce:",2" json:",omitempty"`
	GroupRankSeq      int64 `jce:",3" json:",omitempty"`
	GroupInfoExtraSeq int64 `jce:",4" json:",omitempty"`
}

type FriendListGetGroupListResponse struct {
	Uin                int64           `jce:",0" json:",omitempty"`
	GroupCount         int16           `jce:",1" json:",omitempty"`
	Result             int32           `jce:",2" json:",omitempty"`
	ErrorCode          int16           `jce:",3" json:",omitempty"`
	Cookie             []byte          `jce:",4" json:",omitempty"`
	GroupList          []GroupInfo     `jce:",5" json:",omitempty"`
	GroupListDelete    []GroupInfo     `jce:",6" json:",omitempty"`
	GroupRankList      []GroupRankInfo `jce:",7" json:",omitempty"`
	FavouriteGroupList []FavoriteGroup `jce:",8" json:",omitempty"`
	GroupListExtra     []GroupInfo     `jce:",9" json:",omitempty"`
	GroupInfoExtra     []int64         `jce:",10" json:",omitempty"`
}

type FavoriteGroup struct {
	GroupCode     int64 `jce:",0" json:",omitempty"`
	Timestamp     int64 `jce:",1" json:",omitempty"`
	SNSFlag       int64 `jce:",2" json:",omitempty"`
	OpenTimestamp int64 `jce:",3" json:",omitempty"`
}

type GroupRankInfo struct {
	GroupCode            int64           `jce:",0" json:",omitempty"`
	GroupRankSysFlag     int8            `jce:",1" json:",omitempty"`
	GroupRankUserFlag    int8            `jce:",2" json:",omitempty"`
	RankMap              []LevelRankPair `jce:",3" json:",omitempty"`
	GroupRankSeq         int64           `jce:",4" json:",omitempty"`
	OwnerName            string          `jce:",5" json:",omitempty"`
	AdminName            string          `jce:",6" json:",omitempty"`
	OfficeMode           int64           `jce:",7" json:",omitempty"`
	GroupRankUserFlagNew int8            `jce:",8" json:",omitempty"`
	RankMapNew           []LevelRankPair `jce:",9" json:",omitempty"`
}

type LevelRankPair struct {
	Level int64  `jce:",0" json:",omitempty"`
	Rank  string `jce:",1" json:",omitempty"`
}

type GroupInfo struct {
	GroupUin              int64  `jce:",0" json:",omitempty"`
	GroupCode             int64  `jce:",1" json:",omitempty"`
	Flag                  int8   `jce:",2" json:",omitempty"`
	GroupInfoSeq          int64  `jce:",3" json:",omitempty"`
	GroupName             string `jce:",4" json:",omitempty"`
	GroupMemo             string `jce:",5" json:",omitempty"`
	GroupFlagExt          int64  `jce:",6" json:",omitempty"`
	GroupRankSeq          int64  `jce:",7" json:",omitempty"`
	CertificationType     int64  `jce:",8" json:",omitempty"`
	ShutupTimestamp       int64  `jce:",9" json:",omitempty"`
	MyShutupTimestamp     int64  `jce:",10" json:",omitempty"`
	CmdUinUinFlag         int64  `jce:",11" json:",omitempty"`
	AdditionalFlag        int64  `jce:",12" json:",omitempty"`
	GroupTypeFlag         int64  `jce:",13" json:",omitempty"`
	GroupSecType          int64  `jce:",14" json:",omitempty"`
	GroupSecTypeInfo      int64  `jce:",15" json:",omitempty"`
	GroupClassExt         int64  `jce:",16" json:",omitempty"`
	AppPrivilegeFlag      int64  `jce:",17" json:",omitempty"`
	SubscriptionUin       int64  `jce:",18" json:",omitempty"`
	MemberNum             int64  `jce:",19" json:",omitempty"`
	MemberNumSeq          int64  `jce:",20" json:",omitempty"`
	MemberCardSeq         int64  `jce:",21" json:",omitempty"`
	GroupFlagExt3         int64  `jce:",22" json:",omitempty"`
	GroupOwnerUin         int64  `jce:",23" json:",omitempty"`
	IsConfGroup           bool   `jce:",24" json:",omitempty"`
	IsModifyConfGroupFace bool   `jce:",25" json:",omitempty"`
	IsModifyConfGroupName bool   `jce:",26" json:",omitempty"`
	CmduinJoinTime        int64  `jce:",27" json:",omitempty"`
	CompanyID             int64  `jce:",28" json:",omitempty"`
	MaxGroupMemberNum     int64  `jce:",29" json:",omitempty"`
	CmdUinGroupMask       int64  `jce:",30" json:",omitempty"`
	HLGuildAppid          int64  `jce:",31" json:",omitempty"`
	HLGuildSubType        int64  `jce:",32" json:",omitempty"`
	CmdUinRingtoneID      int64  `jce:",33" json:",omitempty"`
	CmdUinFlagEx2         int64  `jce:",34" json:",omitempty"`
	GroupFlagExt4         int64  `jce:",35" json:",omitempty"`
	AppealDeadline        int64  `jce:",36" json:",omitempty"`
	GroupFlag             int64  `jce:",37" json:",omitempty"`
	GroupRemark           []byte `jce:",38" json:",omitempty"`
}

func NewFriendListGetGroupListRequest(
	uin int64,
	cookie []byte,
) *FriendListGetGroupListRequest {
	return &FriendListGetGroupListRequest{
		Uin:               uin,
		GetMSFMessageFlag: 0x00,
		Cookie:            cookie,
		GroupInfoList:     nil,
		GroupFlagExtra:    0x01,
		Version:           0x00000009,
		CompanyID:         0x0000000000000000,
		VersionNumber:     0x0000000000000001,
		GetLongGroupName:  0x01,
	}
}

func (c *Client) FriendListGetGroupList(
	ctx context.Context,
	req *FriendListGetGroupListRequest,
) (*FriendListGetGroupListResponse, error) {
	buf, err := uni.Marshal(ctx, &uni.Message{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   c.getNextRequestSeq(),
		ServantName: "mqq.IMService.FriendListServiceServantObj",
		FuncName:    "GetTroopListReqV2Simplify",
		Buffer:      []byte{},
		Timeout:     0x00000000,
		Context:     map[string]string{},
		Status:      map[string]string{},
	}, map[string]interface{}{
		"GetTroopListReqV2Simplify": req,
	})
	if err != nil {
		return nil, err
	}
	c2s, s2c := codec.ClientToServerMessage{
		Username: strconv.FormatInt(int64(req.Uin), 10),
		Buffer:   buf,
		Simple:   true,
	}, codec.ServerToClientMessage{}
	err = c.rpc.Call(ServiceMethodFriendListGetGroupList, &c2s, &s2c)
	if err != nil {
		return nil, err
	}
	msg := uni.Message{}
	resp := FriendListGetGroupListResponse{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"GetTroopListRespV2": &resp,
	}); err != nil {
		return nil, err
	}

	dumpServerToClientMessage(&s2c, &resp)

	uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
	for _, item := range resp.GroupList {
		channel := db.Channel{
			Uin:       item.GroupUin,
			Seq:       int32(item.GroupInfoSeq),
			Name:      item.GroupName,
			Memo:      item.GroupMemo,
			MemberNum: int32(item.MemberNum),
			MemberSeq: int32(item.MemberNumSeq),
		}
		c.channels[item.GroupUin] = &channel
		if c.db != nil {
			err := c.dbInsertChannel(uin, &channel)
			if err != nil {
				log.Error().Err(err).Msg(">>> [db  ] dbInsertChannel")
			}
		}
		subResp := &FriendListGetGroupMemberListResponse{}
		for {
			subResp, err = c.FriendListGetGroupMemberList(ctx, NewFriendListGetGroupMemberListRequest(
				req.Uin, item.GroupUin, subResp.NextUin, item.GroupUin,
			))
			if err != nil {
				return nil, err
			}
			if subResp.NextUin == 0 {
				break
			}
		}
	}
	return &resp, nil
}
