package client

import "strconv"

type AuthGetSessionTicketsRequest interface {
	// GetTLVs(ctx context.Context) (map[uint16]tlv.TLVCodec, error)
	GetType() uint16
	SetType(typ uint16)
	GetSeq() uint32
	SetSeq(seq uint32)
	GetUin() uint64
	GetUsername() string
	SetUsername(username string)
	GetServiceMethod() string
}

type AuthGetSessionTicketsResponse struct {
	Code     uint8
	Uin      uint64
	Username string

	PictureSign  []byte
	PictureData  []byte
	CaptchaSign  string
	ErrorCode    uint32
	ErrorTitle   string
	ErrorMessage string
	Message      string
	SMSMobile    string

	AuthSession []byte
	T119        []byte
	T150        []byte
	T161        []byte
	T174        []byte
	T17B        []byte
	T401        [16]byte
	T402        []byte
	T403        []byte
	T546        []byte

	LoginExtraData []byte
}

type authGetSessionTicketsRequest struct {
	typ      uint16
	seq      uint32
	uin      uint64
	username string
	method   string
}

func (req *authGetSessionTicketsRequest) GetType() uint16 {
	return req.typ
}

func (req *authGetSessionTicketsRequest) SetType(typ uint16) {
	req.typ = typ
}

func (req *authGetSessionTicketsRequest) GetSeq() uint32 {
	return req.seq
}

func (req *authGetSessionTicketsRequest) SetSeq(seq uint32) {
	req.seq = seq
}

func (req *authGetSessionTicketsRequest) GetUin() uint64 {
	return req.uin
}

func (req *authGetSessionTicketsRequest) GetUsername() string {
	return req.username
}

func (req *authGetSessionTicketsRequest) SetUsername(username string) {
	req.username = username
	uin, _ := strconv.ParseInt(username, 10, 64)
	req.uin = uint64(uin)
}

func (req *authGetSessionTicketsRequest) GetServiceMethod() string {
	return req.method
}

func (req *authGetSessionTicketsRequest) SetServiceMethod(method string) {
	req.method = method
}
