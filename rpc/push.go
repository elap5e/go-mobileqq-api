package rpc

type PushServiceRequest struct {
	Uin uint64 `jce:",0" json:",omitempty"`
	Map uint64 `jce:",1" json:",omitempty"`
	Str uint8  `jce:",2" json:",omitempty"`
}

type AppPushInfo struct {
	A                   uint32                `jce:",1" json:",omitempty"`
	B                   string                `jce:",2" json:",omitempty"`
	Bid                 uint64                `jce:",3" json:",omitempty"`
	D                   uint64                `jce:",4" json:",omitempty"`
	E                   uint64                `jce:",5" json:",omitempty"`
	F                   uint64                `jce:",6" json:",omitempty"`
	G                   uint64                `jce:",7" json:",omitempty"`
	H                   uint64                `jce:",8" json:",omitempty"`
	I                   string                `jce:",9" json:",omitempty"`
	J                   string                `jce:",10" json:",omitempty"`
	AccountUpdateStatus AccountUpdateStatus   `jce:",10" json:",omitempty"`
	L                   NotifyRegisterInfo    `jce:",11" json:",omitempty"`
	M                   CommandCallbackerInfo `jce:",12" json:",omitempty"`
	N                   string                `jce:",13" json:",omitempty"`
}

type NotifyRegisterInfo struct{}

type CommandCallbackerInfo struct{}
