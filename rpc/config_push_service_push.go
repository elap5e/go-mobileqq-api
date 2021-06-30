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

type SSOServerList struct {
	TGTGList     []SSOServerListInfo `jce:",1"`
	WiFiList     []SSOServerListInfo `jce:",3"`
	Reconnect    uint32              `jce:",4"`
	TestSpeed    bool                `jce:",5"`
	UseNewList   bool                `jce:",6"`
	MultiConn    uint32              `jce:",7"`
	HTTP2G3GList []SSOServerListInfo `jce:",8"`
	HTTPWiFiList []SSOServerListInfo `jce:",9"`
}

type SSOServerListInfo struct {
	IP           string `jce:",1"`
	Port         uint32 `jce:",2"`
	LinkType     bool   `jce:",3"`
	Proxy        bool   `jce:",4"`
	ProtocolType bool   `jce:",5"`
	TimeOut      uint32 `jce:",6"`
}

type FileStorageServerList struct {
	UpLoadList           []FileStorageServerListInfo `jce:",0"`
	PicDownLoadList      []FileStorageServerListInfo `jce:",1"`
	GPicDownLoadList     []FileStorageServerListInfo `jce:",2"`
	QzoneProxyServerList []FileStorageServerListInfo `jce:",3"`
	UrlEncodeServerList  []FileStorageServerListInfo `jce:",4"`
	BigDataIPChannel     BigDataIPChannel            `jce:",5"`
	VIPEmotionList       []FileStorageServerListInfo `jce:",6"`
	C2CPicDownList       []FileStorageServerListInfo `jce:",7"`
	FormatIPInfo         FormatIPInfo                `jce:",8"`
	DomainIPChannel      DomainIPChannel             `jce:",9"`
}

type FileStorageServerListInfo struct {
	IP   string `jce:",1"`
	Port uint32 `jce:",2"`
}

type FormatIPInfo struct {
	IP       string `jce:",0"`
	Operator uint64 `jce:",1"`
}

type BigDataIPChannel struct {
	IPLists []BigDataIPList `jce:",0"`
	Sig     []byte          `jce:",1"`
	Key     []byte          `jce:",2"`
	Uin     uint64          `jce:",3"`
	Flag    uint32          `jce:",4"`
}

type BigDataIPList struct {
	Type    uint64          `jce:",0"`
	IPList  []BigDataIPInfo `jce:",1"`
	Configs []NetSegConf    `jce:",2"`
	Size    uint64          `jce:",3"`
}

type BigDataIPInfo struct {
	Type uint64 `jce:",0"`
	IP   string `jce:",1"`
	Port uint64 `jce:",2"`
}

type NetSegConf struct {
	NetType    uint64 `jce:",0"`
	SegSize    uint64 `jce:",1"`
	SegNum     uint64 `jce:",2"`
	CurConnNum uint64 `jce:",3"`
}

type DomainIPChannel struct {
	IPLists []DomainIPList `jce:",0"`
}

type DomainIPList struct {
	Type   uint32         `jce:",0"`
	IPList []DomainIPInfo `jce:",1"`
}

type DomainIPInfo struct {
	IP   uint32 `jce:",1"`
	Port uint32 `jce:",2"`
}

type ClientLogConfig struct {
	Type       uint32    `jce:",1"`
	TimeStart  TimeStamp `jce:",2"`
	TimeFinish TimeStamp `jce:",3"`
	LogLevel   uint8     `jce:",4"`
	Cookie     uint32    `jce:",5"`
	Seq        uint64    `jce:",6"`
}

type TimeStamp struct {
	Year  uint32 `jce:",1"`
	Month uint8  `jce:",2"`
	Day   uint8  `jce:",3"`
	Hour  uint8  `jce:",4"`
}

type ProxyIPChannel struct {
	ProxyIPLists []ProxyIPList `jce:",0"`
	Reconnect    uint32        `jce:",1"`
}

type ProxyIPList struct {
	Type   uint64        `jce:",0"`
	IPlist []ProxyIPInfo `jce:",1"`
}

type ProxyIPInfo struct {
	Type uint32 `jce:",0"`
	IP   uint32 `jce:",1"`
	Port uint32 `jce:",2"`
}

func (c *Client) handleConfigPushServicePush(
	ctx context.Context,
	s2c *ServerToClientMessage,
) error {
	msg := new(uni.Message)
	req := new(ConfigPushServicePushRequest)
	if err := uni.Unmarshal(ctx, s2c.Buffer, msg, map[string]interface{}{
		"PushReq": msg,
	}); err != nil {
		return err
	}
	resp := new(ConfigPushServicePushResponse)
	switch req.Type {
	case 0x01:
		data := new(SSOServerList)
		if err := jce.Unmarshal(req.Buffer, data); err != nil {
			return err
		}
		// TODO: process message
		resp = &ConfigPushServicePushResponse{
			Type:   req.Type,
			Seq:    req.Seq,
			Buffer: req.Buffer,
		}
	case 0x02:
		data := new(FileStorageServerList)
		if err := jce.Unmarshal(req.Buffer, data); err != nil {
			return err
		}
		// TODO: process message
		resp = &ConfigPushServicePushResponse{
			Type:   req.Type,
			Seq:    req.Seq,
			Buffer: req.Buffer,
		}
	case 0x03:
		data := new(ClientLogConfig)
		if err := jce.Unmarshal(req.Buffer, data); err != nil {
			return err
		}
		// TODO: process message
		resp = &ConfigPushServicePushResponse{
			Type:   req.Type,
			Seq:    req.Seq,
			Buffer: req.Buffer,
		}
	case 0x04:
		data := new(ProxyIPChannel)
		if err := jce.Unmarshal(req.Buffer, data); err != nil {
			return err
		}
		// TODO: process message
		resp = &ConfigPushServicePushResponse{
			Type:   req.Type,
			Seq:    req.Seq,
			Buffer: req.Buffer,
		}
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
	if err := c.Call(ServiceMethodConfigPushServicePushResponse,
		&ClientToServerMessage{
			Username: s2c.Username,
			Seq:      s2c.Seq,
			Buffer:   buf,
			Simple:   false,
		}, s2c,
	); err != nil {
		return err
	}
	return nil
}
