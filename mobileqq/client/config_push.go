package client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/elap5e/go-mobileqq-api/encoding/jce"
	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type ConfigPushRequest struct {
	Type   uint32 `jce:",1" json:"type,omitempty"`
	Seq    uint64 `jce:",3" json:"seq,omitempty"`
	Buffer []byte `jce:",2" json:"buffer,omitempty"`
}

type ConfigPushResponse struct {
	Type   uint32 `jce:",1" json:"type,omitempty"`
	Seq    uint64 `jce:",2" json:"seq,omitempty"`
	Buffer []byte `jce:",3" json:"buffer,omitempty"`
}

type SSOServerConfig struct {
	TGTGList     []SSOServerInfo `jce:",1" json:"2g3g_list,omitempty"`
	WiFiList     []SSOServerInfo `jce:",3" json:"wifi_list,omitempty"`
	Reconnect    uint32          `jce:",4" json:"reconnect,omitempty"`
	TestSpeed    bool            `jce:",5" json:"test_speed,omitempty"`
	UseNewList   bool            `jce:",6" json:"use_new_list,omitempty"`
	MultiConn    uint32          `jce:",7" json:"multi_conn,omitempty"`
	HTTP2G3GList []SSOServerInfo `jce:",8" json:"http_2g3g_list,omitempty"`
	HTTPWiFiList []SSOServerInfo `jce:",9" json:"http_wifi_list,omitempty"`

	Unknown12 []uint64 `jce:",12" json:"unknown12,omitempty"`
	Unknown13 []uint64 `jce:",13" json:"unknown13,omitempty"`
	Unknown14 uint64   `jce:",14" json:"unknown14,omitempty"`
	Unknown15 uint64   `jce:",15" json:"unknown15,omitempty"`
	Unknown16 string   `jce:",16" json:"unknown16,omitempty"`
} // SsoServerList

type SSOServerInfo struct {
	IP           string `jce:",1" json:"ip,omitempty"`
	Port         uint32 `jce:",2" json:"port,omitempty"`
	LinkType     bool   `jce:",3" json:"link_type,omitempty"`
	Proxy        bool   `jce:",4" json:"proxy,omitempty"`
	ProtocolType bool   `jce:",5" json:"protocol_type,omitempty"`
	Timeout      uint32 `jce:",6" json:"timeout,omitempty"`
	Location     string `jce:",8" json:"location,omitempty"`
} // SsoServerListInfo

type FileStorageServerConfig struct {
	UpLoadList               []FileStorageServerInfo `jce:",0" json:"upload_list,omitempty"`
	PictureDownLoadList      []FileStorageServerInfo `jce:",1" json:"picture_downLoad_list,omitempty"`
	GroupPictureDownLoadList []FileStorageServerInfo `jce:",2" json:"group_picture_downLoad_list,omitempty"`
	QZoneProxyServerList     []FileStorageServerInfo `jce:",3" json:"qzone_proxy_server_list,omitempty"`
	URLEncodeServerList      []FileStorageServerInfo `jce:",4" json:"url_encode_server_list,omitempty"`
	BigDataIPChannel         *BigDataIPChannel       `jce:",5" json:"big_data_ip_channel,omitempty"`
	VIPEmotionList           []FileStorageServerInfo `jce:",6" json:"vip_emotion_list,omitempty"`
	C2CPictureDownLoadList   []FileStorageServerInfo `jce:",7" json:"c2c_picture_downLoad_list,omitempty"`
	FormatIPInfo             *FormatIPInfo           `jce:",8" json:"format_ip_info,omitempty"`
	DomainIPChannel          *DomainIPChannel        `jce:",9" json:"domain_ip_channel,omitempty"`
	PTTList                  []byte                  `jce:",10" json:"ptt_list,omitempty"`
}

type FileStorageServerInfo struct {
	IP   string `jce:",1" json:"ip,omitempty"`
	Port uint32 `jce:",2" json:"port,omitempty"`
}

type FormatIPInfo struct {
	IP       string `jce:",0" json:"ip,omitempty"`
	Operator uint64 `jce:",1" json:"operator,omitempty"`
}

type BigDataIPChannel struct {
	BigDataIPList []BigDataIP `jce:",0" json:"big_data_ip_list,omitempty"`
	Sig           []byte      `jce:",1" json:"sig,omitempty"`
	Key           []byte      `jce:",2" json:"key,omitempty"`
	Uin           uint64      `jce:",3" json:"uin,omitempty"`
	Flag          uint32      `jce:",4" json:"flag,omitempty"`
	Buffer        []byte      `jce:",5" json:"buffer,omitempty"`
}

type BigDataIP struct {
	Type       uint64          `jce:",0" json:"type,omitempty"`
	IPList     []BigDataIPInfo `jce:",1" json:"ip_list,omitempty"`
	ConfigList []NetSegConfig  `jce:",2" json:"config_list,omitempty"`
	Size       uint64          `jce:",3" json:"size,omitempty"`
}

type BigDataIPInfo struct {
	Type uint64 `jce:",0" json:"type,omitempty"`
	IP   string `jce:",1" json:"ip,omitempty"`
	Port uint64 `jce:",2" json:"port,omitempty"`
}

type NetSegConfig struct {
	NetType           uint64 `jce:",0" json:"net_type,omitempty"`
	SegSize           uint64 `jce:",1" json:"seg_size,omitempty"`
	SegNumber         uint64 `jce:",2" json:"seg_number,omitempty"`
	CurrentConnNumber uint64 `jce:",3" json:"current_conn_number,omitempty"`
}

type DomainIPChannel struct {
	DomainIPList []DomainIP `jce:",0" json:"domain_ip_list,omitempty"`
}

type DomainIP struct {
	Type   uint32         `jce:",0" json:"type,omitempty"`
	IPList []DomainIPInfo `jce:",1" json:"ip_list,omitempty"`
}

type DomainIPInfo struct {
	IP   uint32 `jce:",1" json:"ip,omitempty"`
	Port uint32 `jce:",2" json:"port,omitempty"`
}

type ClientLogConfig struct {
	Type       uint32     `jce:",1" json:"type,omitempty"`
	TimeStart  *Timestamp `jce:",2" json:"time_start,omitempty"`
	TimeFinish *Timestamp `jce:",3" json:"time_finish,omitempty"`
	LogLevel   uint8      `jce:",4" json:"log_level,omitempty"`
	Cookie     uint32     `jce:",5" json:"cookie,omitempty"`
	Seq        uint64     `jce:",6" json:"seq,omitempty"`
}

type Timestamp struct {
	Year  uint32 `jce:",1" json:"year,omitempty"`
	Month uint8  `jce:",2" json:"month,omitempty"`
	Day   uint8  `jce:",3" json:"day,omitempty"`
	Hour  uint8  `jce:",4" json:"hour,omitempty"`
}

type ProxyIPChannel struct {
	ProxyIPList []ProxyIP `jce:",0" json:"proxy_ip_list,omitempty"`
	Reconnect   uint32    `jce:",1" json:"reconnect,omitempty"`
}

type ProxyIP struct {
	Type   uint64        `jce:",0" json:"type,omitempty"`
	IPlist []ProxyIPInfo `jce:",1" json:"ip_list,omitempty"`
}

type ProxyIPInfo struct {
	Type uint32 `jce:",0" json:"type,omitempty"`
	IP   uint32 `jce:",1" json:"ip,omitempty"`
	Port uint32 `jce:",2" json:"port,omitempty"`
}

func (c *Client) handleConfigPushRequest(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	msg := uni.Message{}
	req := ConfigPushRequest{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"PushReq": &req,
	}); err != nil {
		return nil, err
	}
	switch req.Type {
	case 0x01:
		data := SSOServerConfig{}
		if err := jce.Unmarshal(req.Buffer, &data, true); err != nil {
			return nil, err
		}
		tdata, err := json.MarshalIndent(&data, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(path.Join(
			c.GetCacheByUsernameDir(s2c.Username), "sso-server-config.json",
		), append(tdata, '\n'), 0600); err != nil {
			return nil, err
		}
	case 0x02:
		data := FileStorageServerConfig{}
		if err := jce.Unmarshal(req.Buffer, &data, true); err != nil {
			return nil, err
		}
		tdata, err := json.MarshalIndent(&data, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(path.Join(
			c.GetCacheByUsernameDir(s2c.Username), "file-storage-server-config.json",
		), append(tdata, '\n'), 0600); err != nil {
			return nil, err
		}
	case 0x03:
		data := ClientLogConfig{}
		if err := jce.Unmarshal(req.Buffer, &data, true); err != nil {
			return nil, err
		}
		log.Debug().Msg(">>> [dump] req.Buffer 0x03:\n" + hex.Dump(req.Buffer))
	case 0x04:
		data := ProxyIPChannel{}
		if err := jce.Unmarshal(req.Buffer, &data, true); err != nil {
			return nil, err
		}
		log.Debug().Msg(">>> [dump] req.Buffer 0x04:\n" + hex.Dump(req.Buffer))
	}
	resp := &ConfigPushResponse{
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
		ServiceMethod: ServiceMethodConfigPushResponse,
		Buffer:        buf,
		Simple:        true,
	}, nil
}
