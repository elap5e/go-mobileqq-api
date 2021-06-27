package uni

type Message struct {
	Version     uint16                 `jce:",1"`
	PacketType  uint8                  `jce:",2"`
	MessageType uint32                 `jce:",3"`
	RequestID   uint32                 `jce:",4"`
	ServantName string                 `jce:",5"`
	FuncName    string                 `jce:",6"`
	Buffer      map[string]interface{} `jce:",7"`
	Timeout     uint32                 `jce:",8"`
	Context     map[string]string      `jce:",9"`
	Status      map[string]string      `jce:",10"`
}

type ServicePushRequest struct {
	Uin uint64 `jce:",0"`
	Map uint64 `jce:",1"`
	Str uint8  `jce:",2"`
}

type ServiceRegisterRequest struct {
	Uin          uint64 `jce:",0"`
	Bid          uint64 `jce:",1"`
	ConnType     uint8  `jce:",2"`
	Other        string `jce:",3"`
	Status       uint32 `jce:",4"`
	OnlinePush   bool   `jce:",5"`
	IsOnline     bool   `jce:",6"`
	IsShowOnline bool   `jce:",7"`
	KikPC        bool   `jce:",8"`
	KikWeak      bool   `jce:",9"`
	Timestamp    uint64 `jce:",10"`
	OSVersion    uint32 `jce:",11"`
	NetType      uint8  `jce:",12"`
	BuildVer     string `jce:",13"`
	RegType      bool   `jce:",14"`
	DevParam     []byte `jce:",15"`
	GUID         []byte `jce:",16"`
	LocaleID     uint32 `jce:",17"`
	SlientPush   bool   `jce:",18"`
	DevName      string `jce:",19"`
	DevType      string `jce:",20"`
	OSVer        string `jce:",21"`
	OpenPush     bool   `jce:",22"`
	LargeSeq     uint32 `jce:",23"`
}

func NewServiceRegisterRequest() *Message {
	return &Message{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   0x00000000,
		ServantName: "PushService",
		FuncName:    "SvcReqRegister",
		Buffer: map[string]interface{}{"SvcReqRegister": ServiceRegisterRequest{
			Uin:          0x0000000000000000,
			Bid:          0x0000000000000000,
			ConnType:     0x00,
			Other:        "",
			Status:       0x0000000b,
			OnlinePush:   false,
			IsOnline:     false,
			IsShowOnline: false,
			KikPC:        false,
			KikWeak:      false,
			Timestamp:    0x0000000000000000,
			OSVersion:    0x00000000,
			NetType:      0x00,
			BuildVer:     "",
			RegType:      false,
			DevParam:     nil,
			GUID:         nil,
			LocaleID:     0x00000804,
			SlientPush:   false,
			DevName:      "",
			DevType:      "",
			OSVer:        "",
			OpenPush:     true,
			LargeSeq:     0x00000000,
		}},
		Timeout: 0x00000000,
		Context: nil,
		Status:  nil,
	}
}
