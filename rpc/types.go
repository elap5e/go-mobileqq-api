package rpc

type ClientAuthData struct {
	A1     []byte
	A2     []byte
	A3     []byte
	D1     []byte
	D2     []byte
	S1     []byte
	Key    [16]byte
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
	Seq           uint32
	Username      string
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
	Seq           uint32
	Username      string
	ReturnCode    uint32
	ServiceMethod string
	Cookie        []byte
	Buffer        []byte
}
