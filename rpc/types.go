package rpc

type UserSignature struct {
	Username    string
	DeviceToken []byte `json:",omitempty"`
	Domains     map[string]string
	Tickets     map[string]Ticket
}

type Ticket struct {
	Sig []byte
	Key []byte `json:",omitempty"`
	Iss int64
	Exp int64
}

type ClientCodecKey struct {
	A1     []byte   // tgtgt
	A1Key  []byte   // tgtgtKey
	A2     []byte   // tgt
	A2Key  [16]byte // tgtkey
	A3     []byte
	D1     []byte
	D2     []byte
	D2Key  [16]byte
	S1     []byte
	Cookie []byte
}

type ClientCall struct {
	ServiceMethod         string
	ClientToServerMessage *ClientToServerMessage
	ServerToClientMessage *ServerToClientMessage
	Error                 error
	Done                  chan *ClientCall
}

type ClientToServerMessage struct {
	Version       uint32
	EncryptType   uint8
	EncryptKey    [16]byte
	Username      string
	Seq           uint32
	AppID         uint32
	ServiceMethod string
	Cookie        []byte
	ReserveField  []byte
	Buffer        []byte
	Simple        bool
}

type ServerToClientMessage struct {
	Version       uint32
	EncryptType   uint8
	Username      string
	Seq           uint32
	ReturnCode    uint32
	ServiceMethod string
	Cookie        []byte
	Buffer        []byte
}

type ServerPublicKey struct {
	QuerySpan         uint32 `json:"QuerySpan"`
	PublicKeyMetaData struct {
		KeyVersion    uint16 `json:"KeyVer"`
		PublicKey     string `json:"PubKey"`
		PublicKeySign string `json:"PubKeySign"`
	} `json:"PubKeyMeta"`
}
