package client

import (
	"context"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func NewAccountSetStatusFromClientRequest(
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

func (c *Client) NewAccountSetStatusFromClient(
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
	err = c.rpc.Call(ServiceMethodAccountSetStatusFromClient, &c2s, &s2c)
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
	return &resp, nil
}
