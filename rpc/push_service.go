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
