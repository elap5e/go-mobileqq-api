package rpc

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthCheckSMSAndGetSessionTicketsRequest struct {
	authGetSessionTicketsRequest

	_AuthSession []byte // c.GetUserSignature(req.Username).Session.Auth
	Code         []byte
	_T174        []byte // c.t174
	_MiscBitmap  uint32 // c.cfg.Client.MiscBitmap
	SubSigMap    uint32
	SubAppIDList []uint64
	_HashedGUID  [16]byte // c.hashedGUID
	_ExtraData   []byte   // c.extraData[0x0542]

	lockType uint8
}

func NewAuthCheckSMSAndGetSessionTicketsRequest(
	username string,
	code []byte,
) *AuthCheckSMSAndGetSessionTicketsRequest {
	req := &AuthCheckSMSAndGetSessionTicketsRequest{
		_AuthSession: nil,
		Code:         code,
		_T174:        nil,
		_MiscBitmap:  0x00000000,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		_HashedGUID:  [16]byte{},
		_ExtraData:   nil,

		lockType: 0x00,
	}
	req.SetUsername(username)
	return req
}

func (req *AuthCheckSMSAndGetSessionTicketsRequest) GetTLVs(
	ctx context.Context,
) (map[uint16]tlv.TLVCodec, error) {
	c := ForClient(ctx)
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(
		c.GetUserSignature(req.GetUsername()).Session.Auth,
	)
	tlvs[0x0116] = tlv.NewT116(
		c.cfg.Client.MiscBitmap,
		req.SubSigMap,
		req.SubAppIDList,
	)
	tlvs[0x0174] = tlv.NewT174(c.t174)
	tlvs[0x017c] = tlv.NewT17C(req.Code)
	tlvs[0x0401] = tlv.NewT401(c.hashedGUID)
	tlvs[0x0197] = tlv.NewTLV(
		0x0198,
		0x0000,
		bytes.NewBuffer([]byte{req.lockType}),
	)
	tlvs[0x0542] = tlv.NewT542(c.extraData[0x0542])
	req.SetType(0x0007)
	req.SetServiceMethod(ServiceMethodAuthLogin)
	return tlvs, nil
}

func (c *Client) AuthCheckSMSAndGetSessionTickets(
	ctx context.Context,
	req *AuthCheckSMSAndGetSessionTicketsRequest,
) (*AuthGetSessionTicketsResponse, error) {
	return c.AuthGetSessionTickets(ctx, req)
}
