package client

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthCheckCaptchaAndGetSessionTicketsRequest struct {
	authGetSessionTicketsRequest

	_AuthSession []byte // c.GetUserSignature(req.Username).Session.Auth
	Code         []byte
	Sign         []byte
	_MiscBitmap  uint32 // c.cfg.Client.MiscBitmap
	SubSigMap    uint32
	SubAppIDList []uint64
	_ExtraData   []byte // c.t547

	isCaptcha bool
}

func NewAuthCheckCaptchaAndGetSessionTicketsRequest(
	username string,
	code []byte,
) *AuthCheckCaptchaAndGetSessionTicketsRequest {
	req := &AuthCheckCaptchaAndGetSessionTicketsRequest{
		_AuthSession: nil,
		Code:         code,
		Sign:         nil, // nil
		_MiscBitmap:  0x00000000,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		_ExtraData:   nil,

		isCaptcha: true,
	}
	req.SetUsername(username)
	return req
}

func (req *AuthCheckCaptchaAndGetSessionTicketsRequest) GetTLVs(
	ctx context.Context,
) (map[uint16]tlv.TLVCodec, error) {
	c := ForClient(ctx)
	tlvs := make(map[uint16]tlv.TLVCodec)
	if req.isCaptcha {
		tlvs[0x0193] = tlv.NewT193(req.Code)
	} else {
		tlvs[0x0002] = tlv.NewT2(req.Code, req.Sign)
	}
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(
		c.rpc.GetUserSignature(req.GetUsername()).Session.Auth,
	)
	tlvs[0x0116] = tlv.NewT116(
		c.cfg.Client.MiscBitmap,
		req.SubSigMap,
		req.SubAppIDList,
	)
	tlvs[0x0547] = tlv.NewT547(c.extraData[0x0547])
	req.SetType(0x0002)
	req.SetServiceMethod(ServiceMethodAuthLogin)
	return tlvs, nil
}

func (c *Client) AuthCheckCaptchaAndGetSessionTickets(
	ctx context.Context,
	req *AuthCheckCaptchaAndGetSessionTicketsRequest,
) (*AuthGetSessionTicketsResponse, error) {
	return c.AuthGetSessionTickets(ctx, req)
}