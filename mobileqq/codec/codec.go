package codec

type ClientCodec interface {
	Close() error

	ReadBody(msg *ServerToClientMessage) error
	ReadHead(msg *ServerToClientMessage) error
	Write(msg *ClientToServerMessage) error
}

type ClientToServerMessage struct {
	Simple      bool
	Version     uint32
	EncryptType uint8
	UserD2      []byte
	Username    string

	UserD2Key [16]byte

	Seq           uint32
	FixID         uint32
	AppID         uint32
	NetworkType   uint8 // 0x00: Others; 0x01: Wi-Fi
	NetIPFamily   uint8 // 0x00: Others; 0x01: IPv4; 0x02: IPv6; 0x03: Dual
	UserA2        []byte
	ServiceMethod string
	Cookie        []byte
	IMEI          string
	KSID          []byte
	IMSI          string
	Revision      string
	ReserveField  []byte
	Buffer        []byte
}

type ServerToClientMessage struct {
	Version     uint32
	EncryptType uint8
	Username    string

	UserD2Key [16]byte

	Seq           uint32
	Code          uint32
	Message       string
	ServiceMethod string
	Cookie        []byte
	Flag          uint32
	ReserveField  []byte
	Buffer        []byte
}

const (
	EncryptTypeNotNeedEncrypt = 0x00
	EncryptTypeEncryptByD2Key = 0x01
	EncryptTypeEncryptByZeros = 0x02
)

const (
	FlagNoCompression   = 0x00000000
	FlagZlibCompression = 0x00000001
)

const (
	VersionDefault = 0x0000000a
	VersionSimple  = 0x0000000b
)

func CopyServerToClientMessage(dst, src *ServerToClientMessage) {
	dst.Version = src.Version
	dst.EncryptType = src.EncryptType
	dst.Username = src.Username
	dst.UserD2Key = src.UserD2Key
	dst.Seq = src.Seq
	dst.Code = src.Code
	dst.Message = src.Message
	dst.ServiceMethod = src.ServiceMethod
	dst.Cookie = src.Cookie
	dst.Flag = src.Flag
	dst.ReserveField = src.ReserveField
	dst.Buffer = src.Buffer
}
