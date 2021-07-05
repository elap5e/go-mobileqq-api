package client

import (
	"context"
)

type AuthGetSessionTicketsWithoutPasswordRequest struct {
	authGetSessionTicketsRequest

	_Uin             uint64
	DstAppID         uint64
	SrcAppID         uint64 // constant 0x00000064
	AppClientVersion uint32
	MainSigMap       uint32
	_UserA2          []byte // sig.Tickets["A2"].Sig
	_MiscBitmap      uint32 // c.cfg.Client.MiscBitmap
	SubSigMap        uint32
	SubAppIDList     []uint64
	_KSID            []byte // sig.Session.KSID
	_UserD2          []byte // sig.Tickets["D2"].Sig
	Domains          []string

	changeD2 bool
}

func (c *Client) AuthGetSessionTicketsWithoutPassword(
	ctx context.Context,
	req *AuthGetSessionTicketsWithoutPasswordRequest,
) (*AuthGetSessionTicketsResponse, error) {
	return c.AuthGetSessionTickets(ctx, req)
}
