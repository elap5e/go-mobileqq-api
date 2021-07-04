package rpc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/elap5e/go-mobileqq-api/encoding/jce"
	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type PushConfigRequest struct {
	Type   uint32 `jce:",1" json:",omitempty"`
	Seq    uint64 `jce:",3" json:",omitempty"`
	Buffer []byte `jce:",2" json:",omitempty"`
}

type PushConfigResponse struct {
	Type   uint32 `jce:",1" json:",omitempty"`
	Seq    uint64 `jce:",2" json:",omitempty"`
	Buffer []byte `jce:",3" json:",omitempty"`
}

type SSOServerList struct {
	TGTGList     []SSOServerListInfo `jce:",1" json:",omitempty"`
	WiFiList     []SSOServerListInfo `jce:",3" json:",omitempty"`
	Reconnect    uint32              `jce:",4" json:",omitempty"`
	TestSpeed    bool                `jce:",5" json:",omitempty"`
	UseNewList   bool                `jce:",6" json:",omitempty"`
	MultiConn    uint32              `jce:",7" json:",omitempty"`
	HTTP2G3GList []SSOServerListInfo `jce:",8" json:",omitempty"`
	HTTPWiFiList []SSOServerListInfo `jce:",9" json:",omitempty"`

	Unknown12 []uint64 `jce:",12" json:",omitempty"`
	Unknown13 []uint64 `jce:",13" json:",omitempty"`
	Unknown14 uint64   `jce:",14" json:",omitempty"`
	Unknown15 uint64   `jce:",15" json:",omitempty"`
	Unknown16 string   `jce:",16" json:",omitempty"`
}

type SSOServerListInfo struct {
	IP           string `jce:",1" json:",omitempty"`
	Port         uint32 `jce:",2" json:",omitempty"`
	LinkType     bool   `jce:",3" json:",omitempty"`
	Proxy        bool   `jce:",4" json:",omitempty"`
	ProtocolType bool   `jce:",5" json:",omitempty"`
	Timeout      uint32 `jce:",6" json:",omitempty"`

	Unknown8 string `jce:",8" json:",omitempty"`
}

type FileStorageServerList struct {
	UpLoadList           []FileStorageServerListInfo `jce:",0" json:",omitempty"`
	PicDownLoadList      []FileStorageServerListInfo `jce:",1" json:",omitempty"`
	GPicDownLoadList     []FileStorageServerListInfo `jce:",2" json:",omitempty"`
	QzoneProxyServerList []FileStorageServerListInfo `jce:",3" json:",omitempty"`
	UrlEncodeServerList  []FileStorageServerListInfo `jce:",4" json:",omitempty"`
	BigDataIPChannel     *BigDataIPChannel           `jce:",5" json:",omitempty"`
	VIPEmotionList       []FileStorageServerListInfo `jce:",6" json:",omitempty"`
	C2CPicDownList       []FileStorageServerListInfo `jce:",7" json:",omitempty"`
	FormatIPInfo         *FormatIPInfo               `jce:",8" json:",omitempty"`
	DomainIPChannel      *DomainIPChannel            `jce:",9" json:",omitempty"`
	PTTList              []byte                      `jce:",10" json:",omitempty"`
}

type FileStorageServerListInfo struct {
	IP   string `jce:",1" json:",omitempty"`
	Port uint32 `jce:",2" json:",omitempty"`
}

type FormatIPInfo struct {
	IP       string `jce:",0" json:",omitempty"`
	Operator uint64 `jce:",1" json:",omitempty"`
}

type BigDataIPChannel struct {
	IPLists []BigDataIPList `jce:",0" json:",omitempty"`
	Sig     []byte          `jce:",1" json:",omitempty"`
	Key     []byte          `jce:",2" json:",omitempty"`
	Uin     uint64          `jce:",3" json:",omitempty"`
	Flag    uint32          `jce:",4" json:",omitempty"`
	Buffer  []byte          `jce:",5" json:",omitempty"`
}

type BigDataIPList struct {
	Type    uint64          `jce:",0" json:",omitempty"`
	IPList  []BigDataIPInfo `jce:",1" json:",omitempty"`
	Configs []NetSegConf    `jce:",2" json:",omitempty"`
	Size    uint64          `jce:",3" json:",omitempty"`
}

type BigDataIPInfo struct {
	Type uint64 `jce:",0" json:",omitempty"`
	IP   string `jce:",1" json:",omitempty"`
	Port uint64 `jce:",2" json:",omitempty"`
}

type NetSegConf struct {
	NetType    uint64 `jce:",0" json:",omitempty"`
	SegSize    uint64 `jce:",1" json:",omitempty"`
	SegNum     uint64 `jce:",2" json:",omitempty"`
	CurConnNum uint64 `jce:",3" json:",omitempty"`
}

type DomainIPChannel struct {
	IPLists []DomainIPList `jce:",0" json:",omitempty"`
}

type DomainIPList struct {
	Type   uint32         `jce:",0" json:",omitempty"`
	IPList []DomainIPInfo `jce:",1" json:",omitempty"`
}

type DomainIPInfo struct {
	IP   uint32 `jce:",1" json:",omitempty"`
	Port uint32 `jce:",2" json:",omitempty"`
}

type ClientLogConfig struct {
	Type       uint32     `jce:",1" json:",omitempty"`
	TimeStart  *TimeStamp `jce:",2" json:",omitempty"`
	TimeFinish *TimeStamp `jce:",3" json:",omitempty"`
	LogLevel   uint8      `jce:",4" json:",omitempty"`
	Cookie     uint32     `jce:",5" json:",omitempty"`
	Seq        uint64     `jce:",6" json:",omitempty"`
}

type TimeStamp struct {
	Year  uint32 `jce:",1" json:",omitempty"`
	Month uint8  `jce:",2" json:",omitempty"`
	Day   uint8  `jce:",3" json:",omitempty"`
	Hour  uint8  `jce:",4" json:",omitempty"`
}

type ProxyIPChannel struct {
	ProxyIPLists []ProxyIPList `jce:",0" json:",omitempty"`
	Reconnect    uint32        `jce:",1" json:",omitempty"`
}

type ProxyIPList struct {
	Type   uint64        `jce:",0" json:",omitempty"`
	IPlist []ProxyIPInfo `jce:",1" json:",omitempty"`
}

type ProxyIPInfo struct {
	Type uint32 `jce:",0" json:",omitempty"`
	IP   uint32 `jce:",1" json:",omitempty"`
	Port uint32 `jce:",2" json:",omitempty"`
}

func (c *Client) handlePushConfigRequest(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	msg := uni.Message{}
	req := PushConfigRequest{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"PushReq": &req,
	}); err != nil {
		return nil, err
	}
	switch req.Type {
	case 0x01:
		data := SSOServerList{}
		if err := jce.Unmarshal(req.Buffer, &data, true); err != nil {
			return nil, err
		}
		tdata, err := json.MarshalIndent(&data, "", "  ")
		if err != nil {
			return nil, err
		}
		if ioutil.WriteFile(path.Join(
			c.cfg.CacheDir, s2c.Username, "sso-server-list.json",
		), append(tdata, '\n'), 0600); err != nil {
			return nil, err
		}
	case 0x02:
		data := FileStorageServerList{}
		if err := jce.Unmarshal(req.Buffer, &data, true); err != nil {
			return nil, err
		}
		tdata, err := json.MarshalIndent(&data, "", "  ")
		if err != nil {
			return nil, err
		}
		if ioutil.WriteFile(path.Join(
			c.cfg.CacheDir, s2c.Username, "file-storage-server-list.json",
		), append(tdata, '\n'), 0600); err != nil {
			return nil, err
		}
	case 0x03:
		data := ClientLogConfig{}
		if err := jce.Unmarshal(req.Buffer, &data, true); err != nil {
			return nil, err
		}
	case 0x04:
		data := ProxyIPChannel{}
		if err := jce.Unmarshal(req.Buffer, &data, true); err != nil {
			return nil, err
		}
	}
	resp := &PushConfigResponse{
		Type:   req.Type,
		Seq:    req.Seq,
		Buffer: req.Buffer,
	}
	buf, err := uni.Marshal(ctx, &uni.Message{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   0x00000000,
		ServantName: "QQService.PushConfigSvc.MainServant",
		FuncName:    "PushResp",
		Buffer:      []byte{},
		Timeout:     0x00000000,
		Context:     map[string]string{},
		Status:      map[string]string{},
	}, map[string]interface{}{
		"PushResp": resp,
	})
	if err != nil {
		return nil, err
	}
	return &codec.ClientToServerMessage{
		Username:      s2c.Username,
		Seq:           s2c.Seq,
		ServiceMethod: ServiceMethodPushConfigResponse,
		Buffer:        buf,
		Simple:        true,
	}, nil
}
