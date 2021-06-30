package rpc

type PushServiceRequest struct {
	Uin uint64 `jce:",0"`
	Map uint64 `jce:",1"`
	Str uint8  `jce:",2"`
}

type AppPushInfo struct {
	A                   uint32                `jce:",1"`
	B                   string                `jce:",2"`
	Bid                 uint64                `jce:",3"`
	D                   uint64                `jce:",4"`
	E                   uint64                `jce:",5"`
	F                   uint64                `jce:",6"`
	G                   uint64                `jce:",7"`
	H                   uint64                `jce:",8"`
	I                   string                `jce:",9"`
	J                   string                `jce:",10"`
	AccountUpdateStatus AccountUpdateStatus   `jce:",10"`
	L                   NotifyRegisterInfo    `jce:",11"`
	M                   CommandCallbackerInfo `jce:",12"`
	N                   string                `jce:",13"`
}

type NotifyRegisterInfo struct{}

type CommandCallbackerInfo struct{}

type SSOServerList struct {
	TGTGList     []SSOServerListInfo `jce:",1"`
	WifiList     []SSOServerListInfo `jce:",3"`
	Reconnect    uint32              `jce:",4"`
	TestSpeed    bool                `jce:",5"`
	UseNewList   bool                `jce:",6"`
	MultiConn    uint32              `jce:",7"`
	HTTP2G3GList []SSOServerListInfo `jce:",8"`
	HTTPWifiList []SSOServerListInfo `jce:",9"`
}

type SSOServerListInfo struct {
	IP           string `jce:",1"`
	Port         uint32 `jce:",2"`
	LinkType     bool   `jce:",3"`
	Proxy        bool   `jce:",4"`
	ProtocolType bool   `jce:",5"`
	TimeOut      uint32 `jce:",6"`
}

type FileStoragePushFSServiceList struct {
	UpLoadList            []FileStorageServerListInfo `jce:",0"`
	PicDownLoadList       []FileStorageServerListInfo `jce:",1"`
	GPicDownLoadList      []FileStorageServerListInfo `jce:",2"`
	QzoneProxyServiceList []FileStorageServerListInfo `jce:",3"`
	UrlEncodeServiceList  []FileStorageServerListInfo `jce:",4"`
	DigDataChannel        DigDataChannel              `jce:",5"`
	VIPEmotionList        []FileStorageServerListInfo `jce:",6"`
	C2CPicDownList        []FileStorageServerListInfo `jce:",7"`
	FormatIPInfo          FormatIPInfo                `jce:",8"`
	DomainIpChannel       DomainIpChannel             `jce:",9"`
}

type FileStorageServerListInfo struct {
	IP   string `jce:",1"`
	Port uint32 `jce:",2"`
}

type FormatIPInfo struct {
	GateIP         string `jce:",0"`
	GateIPOperator uint64 `jce:",1"`
}

type DigDataChannel struct {
	IPLists     []BigDataIPList `jce:",0"`
	SigSession  []byte          `jce:",1"`
	KeySession  []byte          `jce:",2"`
	SigUin      uint64          `jce:",3"`
	ConnectFlag uint32          `jce:",4"`
}

type BigDataIPList struct {
	ServiceType  uint64          `jce:",0"`
	IPList       []BigDataIPInfo `jce:",1"`
	NetSegConfs  []NetSegConf    `jce:",2"`
	FragmentSize uint64          `jce:",3"`
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

type DomainIpChannel struct {
	DomainIPLists []DomainIpList `jce:",0"`
}

type DomainIpList struct {
	DomainType uint32         `jce:",0"`
	IPList     []DomainIPInfo `jce:",1"`
}

type DomainIPInfo struct {
	IP   uint32 `jce:",1"`
	Port uint32 `jce:",2"`
}
