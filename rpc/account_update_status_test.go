package rpc

import (
	"context"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
)

func TestMarshalAccountUpdateStatus(t *testing.T) {
	msg := &uni.Message{
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
	}
	ids := []uint64{0x01, 0x02, 0x04}
	bid := uint64(0x0000000000000000)
	for _, id := range ids {
		bid |= id
	}
	push := &AppPushInfo{
		Bid: bid,
		AccountUpdateStatus: AccountUpdateStatus{
			Uin:       0x0000000000002710,
			PushIDs:   ids,
			Status:    uint32(PushRegisterInfoStatusOnline),
			KikPC:     false,
			KikWeak:   false,
			Timestamp: 0x0000000000000000,
			LargeSeq:  0x00000000,
		},
	}
	req := &AccountUpdateStatusRequest{
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
		GUID:         nil, // placeholder
		LocaleID:     0x00000804,
		SlientPush:   false,
		DeviceName:   defaultDeviceOSBuildModel,
		DeviceType:   defaultDeviceOSBuildModel,
		OSVersion:    defaultDeviceOSVersion,
		OpenPush:     true,
		LargeSeq:     push.AccountUpdateStatus.LargeSeq,
	}
	type args struct {
		ctx  context.Context
		msg  *uni.Message
		opts map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "0x00",
			args: args{
				ctx: context.Background(),
				msg: msg,
				opts: map[string]interface{}{
					"SvcReqRegister": req,
				},
			},
			want: []byte{
				0x00, 0x00, 0x00, 0x8f, 0x10, 0x03, 0x2c, 0x3c, 0x4c, 0x56, 0x0b, 0x50, 0x75, 0x73, 0x68, 0x53,
				0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x66, 0x0e, 0x53, 0x76, 0x63, 0x52, 0x65, 0x71, 0x52, 0x65,
				0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x7d, 0x00, 0x00, 0x60, 0x08, 0x00, 0x01, 0x06, 0x0e, 0x53,
				0x76, 0x63, 0x52, 0x65, 0x71, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x1d, 0x00, 0x00,
				0x49, 0x0a, 0x01, 0x27, 0x10, 0x10, 0x07, 0x2c, 0x36, 0x00, 0x40, 0x0b, 0x5c, 0x6c, 0x7c, 0x8c,
				0x9c, 0xac, 0xb0, 0x1e, 0xc0, 0x01, 0xd6, 0x00, 0xec, 0xfd, 0x0f, 0x00, 0x0c, 0xfd, 0x10, 0x00,
				0x0c, 0xf1, 0x11, 0x08, 0x04, 0xfc, 0x12, 0xf6, 0x13, 0x09, 0x52, 0x65, 0x64, 0x6d, 0x69, 0x20,
				0x4b, 0x32, 0x30, 0xf6, 0x14, 0x09, 0x52, 0x65, 0x64, 0x6d, 0x69, 0x20, 0x4b, 0x32, 0x30, 0xf6,
				0x15, 0x02, 0x31, 0x31, 0xf0, 0x16, 0x01, 0xfc, 0x17, 0x0b, 0x8c, 0x98, 0x0c, 0xa8, 0x0c,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uni.Marshal(tt.args.ctx, tt.args.msg, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() =\n%s, want\n%s", hex.Dump(got), hex.Dump(tt.want))
			}
		})
	}
}
