package rpc

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthUnlockDeviceRequest struct {
	authGetSessionTicketsRequest

	_AuthSession []byte // c.GetUserSignature(req.Username).Session.Auth
	_MiscBitmap  uint32 // c.cfg.Client.MiscBitmap
	SubSigMap    uint32
	SubAppIDList []uint64
	_HashedGUID  [16]byte // c.hashedGUID
}

func NewAuthUnlockDeviceRequest(username string) *AuthUnlockDeviceRequest {
	req := &AuthUnlockDeviceRequest{
		_AuthSession: nil,
		_MiscBitmap:  0x00000000,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		_HashedGUID:  [16]byte{},
	}
	req.SetUsername(username)
	return req
}

func (req *AuthUnlockDeviceRequest) GetTLVs(ctx context.Context) (map[uint16]tlv.TLVCodec, error) {
	c := ForClient(ctx)
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(c.GetUserSignature(req.GetUsername()).Session.Auth)
	tlvs[0x0116] = tlv.NewT116(c.cfg.Client.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0401] = tlv.NewT401(c.hashedGUID)
	req.SetType(0x0014)
	req.SetServiceMethod(ServiceMethodAuthLogin)
	return tlvs, nil
}

func (c *Client) AuthUnlockDevice(ctx context.Context, req *AuthUnlockDeviceRequest) (*AuthGetSessionTicketsResponse, error) {
	return c.AuthGetSessionTickets(ctx, req)
}
