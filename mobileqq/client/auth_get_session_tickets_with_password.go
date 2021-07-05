package client

import (
	"context"
	"net"
)

type AuthGetSessionTicketsWithPasswordRequest struct {
	authGetSessionTicketsRequest

	DstAppID         uint64
	SubDstAppID      uint64
	AppClientVersion uint32 // constant 0x00000000
	_Uin             uint64
	I2               uint16 // constant 0x0000
	_IPv4Address     net.IP // c.cfg.Client.MiscBitmap
	ServerTime       uint32
	PasswordMD5      [16]byte
	_UserA1Key       [16]byte // c.userA1Key
	LoginType        uint32   // 0x00, 0x01, 0x03
	UserA1           []byte
	T16A             []byte
	_MiscBitmap      uint32 // c.cfg.Client.MiscBitmap
	SubSigMap        uint32
	SubAppIDList     []uint64
	MainSigMap       uint32
	SrcAppID         uint64
	I7               uint16 // constant 0x0000
	I8               uint8  // constant 0x00
	I9               uint16 // constant 0x0000
	I10              uint8  // constant 0x01
	_KSID            []byte // sig.Session.KSID
	_AuthSession     []byte // sig.Session.AuthSession
	_PackageName     []byte // []byte(c.cfg.Client.PackageName)
	Domains          []string
}

func (c *Client) AuthGetSessionTicketsWithPassword(
	ctx context.Context,
	req *AuthGetSessionTicketsWithPasswordRequest,
) (*AuthGetSessionTicketsResponse, error) {
	return c.AuthGetSessionTickets(ctx, req)
}
