package rpc

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
)

type PushMessageNotifyRequest struct {
	Uin         uint64 `jce:",0"`
	Type        uint8  `jce:",1"`
	Service     string `jce:",2"`
	Cmd         string `jce:",3"`
	Cookie      []byte `jce:",4"`
	MsgType     uint16 `jce:",5"`
	UserActive  uint32 `jce:",6"`
	GeneralFlag uint32 `jce:",7"`
	BindedUin   uint64 `jce:",8"`

	MsgInfo    *MsgInfo `jce:",9"`
	MsgCtrlBuf string   `jce:",10"`
	ServerBuf  []byte   `jce:",11"`
	PingFlag   uint64   `jce:",12"`
	VRIP       uint16   `jce:",13"`

	Unknown14 []byte                             `jce:",14"`
	Unknown15 *PushMessageNotifyRequestUnknown15 `jce:",15"`
	Unknown16 *PushMessageNotifyRequestUnknown15 `jce:",16"`
	Unknown17 *PushMessageNotifyRequestUnknown17 `jce:",17"`
}

type PushMessageNotifyRequestUnknown15 struct {
	Unknown0 uint64 `jce:",0"`
	Unknown1 uint64 `jce:",1"`
	Unknown2 uint64 `jce:",2"`
	Unknown3 uint64 `jce:",3"`
	Unknown4 uint64 `jce:",4"`
}

type PushMessageNotifyRequestUnknown17 struct {
	Unknown0 string `jce:",0"`
	Unknown1 string `jce:",1"`
	Unknown2 string `jce:",2"`
}

type MsgInfo struct {
	FromUin        uint64       `jce:"0"`
	MsgTime        uint64       `jce:"1"`
	MsgType        uint16       `jce:"2"`
	MsgSeq         uint16       `jce:"3"`
	Msg            string       `jce:"4"`
	RealMsgTime    uint64       `jce:"5"`
	MsgBytes       []byte       `jce:"6"`
	AppShareID     uint64       `jce:"7"`
	MsgCookies     []byte       `jce:"8"`
	AppShareCookie []byte       `jce:"9"`
	MsgUid         uint64       `jce:"10"`
	LastChangeTime uint64       `jce:"11"`
	CPicInfo       []CPicInfo   `jce:"12"`
	ShareData      *ShareData   `jce:"13"`
	FromInstID     uint64       `jce:"14"`
	RemarkOfSender []byte       `jce:"15"`
	FromMobile     string       `jce:"16"`
	FromName       string       `jce:"17"`
	NickName       []string     `jce:"18"`
	TempMsgHead    *TempMsgHead `jce:"19"`
}

type CPicInfo struct {
	Path []byte `jce:"0"`
	Host []byte `jce:"1"`
}

type ShareData struct {
	Pkgname string `jce:"0"`
	Msgtail string `jce:"1"`
	PicURL  string `jce:"2"`
	URL     string `jce:"3"`
}

type TempMsgHead struct {
	C2CType     uint32 `jce:"0"`
	ServiceType uint32 `jce:"1"`
}

func (c *Client) handlePushMessageNotify(
	ctx context.Context,
	s2c *ServerToClientMessage,
) (*ClientToServerMessage, error) {
	msg := new(uni.Message)
	req := new(PushMessageNotifyRequest)
	if err := uni.Unmarshal(ctx, s2c.Buffer[4:], msg, map[string]interface{}{
		"req_PushNotify": req,
	}); err != nil {
		return nil, err
	}
	c2s, err := c.MessageGetMessage(ctx, s2c)
	if err != nil {
		return nil, err
	}
	// TODO: something
	_ = c2s
	return nil, nil
}
