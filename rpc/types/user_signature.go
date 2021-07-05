package types

type UserSignature struct {
	Username    string             `json:"username"`
	PasswordMD5 []byte             `json:"-"`
	DeviceToken []byte             `json:"deviceToken,omitempty"`
	Domains     map[string]string  `json:"domains,omitempty"`
	Tickets     map[string]*Ticket `json:"tickets,omitempty"`
	Session     *Session           `json:"session"`
}

type Ticket struct {
	Sig []byte `json:"sig"`
	Key []byte `json:"key,omitempty"`
	Iss int64  `json:"iss"`
	Exp int64  `json:"exp,omitempty"`
}

type Session struct {
	Auth   []byte `json:"auth,omitempty"`
	Cookie []byte `json:"cookie"`
	KSID   []byte `json:"ksid,omitempty"`
}
