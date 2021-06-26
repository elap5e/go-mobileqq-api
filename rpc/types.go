package rpc

type ClientCodecKey struct {
	A1     []byte
	A2     []byte
	A2Key  [16]byte
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
