package rpc

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/encoding/jce"
	"github.com/elap5e/go-mobileqq-api/encoding/uni"
)

type ConfigPushServicePushRequest struct {
	Type   uint32 `jce:",1"`
	Seq    uint64 `jce:",3"`
	Buffer []byte `jce:",2"`
}

type ConfigPushServicePushResponse struct {
	Type   uint32 `jce:",1"`
	Seq    uint64 `jce:",2"`
	Buffer []byte `jce:",3"`
}

func (c *Client) handleConfigPushServicePush(ctx context.Context, s2c *ServerToClientMessage) error {
	req := new(ConfigPushServicePushRequest)
	if err := jce.Unmarshal(s2c.Buffer, req, true); err != nil {
		return err
	}
	msg := new(uni.Message)
	if err := uni.Unmarshal(ctx, req.Buffer, msg, map[string]interface{}{
		"PushReq": map[string]interface{}{
			"ConfigPush.PushReq": req,
		},
	}); err != nil {
		return err
	}
	// TODO: process message
	resp := ConfigPushServicePushResponse{
		Type:   req.Type,
		Seq:    req.Seq,
		Buffer: req.Buffer,
	}
	buf, err := uni.Marshal(ctx, &uni.Message{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   0x00000000,
		ServantName: "QQService.ConfigPushSvc.MainServant",
		FuncName:    "PushResp",
		Buffer:      map[string][]byte{},
		Timeout:     0x00000000,
		Context:     map[string]string{},
		Status:      map[string]string{},
	}, map[string]interface{}{
		"PushResp": resp,
	})
	if err != nil {
		return err
	}
	if err := c.Call(ServiceMethodAccountUpdateStatus, &ClientToServerMessage{
		Username: s2c.Username, // placeholder
		Seq:      s2c.Seq,      // placeholder
		Buffer:   buf,
		Simple:   false,
	}, s2c); err != nil {
		return err
	}
	return nil
}
