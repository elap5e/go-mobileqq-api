package client

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"reflect"
)

type Uint32IPType uint32

func (v Uint32IPType) String() string {
	ip := net.IP{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(ip, uint32(v))
	return ip.String()
}

func (v Uint32IPType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + v.String() + "\""), nil
}

func (v *Uint32IPType) UnmarshalJSON(p []byte) error {
	if len(p) < 2 || p[0] != '"' || p[len(p)-1] != '"' {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(Uint32IPType(0))}
	}
	ip := net.ParseIP(string(p[1 : len(p)-1])).To4()
	*v = Uint32IPType(uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3]))
	return nil
}

type PushServiceRequest struct {
	Uin uint64 `jce:",0" json:",omitempty"`
	Map uint64 `jce:",1" json:",omitempty"`
	Str uint8  `jce:",2" json:",omitempty"`
}

type AppPushInfo struct {
	A             uint32                `jce:",1" json:",omitempty"`
	B             string                `jce:",2" json:",omitempty"`
	Bid           uint64                `jce:",3" json:",omitempty"`
	D             uint64                `jce:",4" json:",omitempty"`
	E             uint64                `jce:",5" json:",omitempty"`
	F             uint64                `jce:",6" json:",omitempty"`
	G             uint64                `jce:",7" json:",omitempty"`
	H             uint64                `jce:",8" json:",omitempty"`
	I             string                `jce:",9" json:",omitempty"`
	J             string                `jce:",10" json:",omitempty"`
	AccountStatus AccountStatus         `jce:",10" json:",omitempty"`
	L             NotifyRegisterInfo    `jce:",11" json:",omitempty"`
	M             CommandCallbackerInfo `jce:",12" json:",omitempty"`
	N             string                `jce:",13" json:",omitempty"`
}

type NotifyRegisterInfo struct{}

type CommandCallbackerInfo struct{}
