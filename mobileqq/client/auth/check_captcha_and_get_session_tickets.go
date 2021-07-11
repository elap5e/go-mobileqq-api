package auth

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/tlv"
)

type checkCaptchaAndGetSessionTicketsRequest struct {
	request

	_AuthSession []byte // h.GetUserSignature(req.Username).Session.Auth
	Code         []byte
	Sign         []byte
	_MiscBitmap  uint32 // h.opt.Client.MiscBitmap
	SubSigMap    uint32
	SubAppIDList []uint64
	_ExtraData   []byte // h.t547

	isCaptcha bool
}

func newCheckCaptchaAndGetSessionTicketsRequest(
	username string,
	code []byte,
) *checkCaptchaAndGetSessionTicketsRequest {
	req := &checkCaptchaAndGetSessionTicketsRequest{
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

func (req *checkCaptchaAndGetSessionTicketsRequest) MustGetTLVs(
	ctx context.Context,
) map[uint16]tlv.TLVCodec {
	h := ForHandler(ctx)
	tlvs := make(map[uint16]tlv.TLVCodec)
	if req.isCaptcha {
		tlvs[0x0193] = tlv.NewT193(req.Code)
	} else {
		tlvs[0x0002] = tlv.NewT2(req.Code, req.Sign)
	}
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(
		h.client.GetUserSignature(req.GetUsername()).Session.Auth,
	)
	tlvs[0x0116] = tlv.NewT116(
		h.opt.Client.MiscBitmap,
		req.SubSigMap,
		req.SubAppIDList,
	)
	tlvs[0x0547] = tlv.NewT547(h.extraData[0x0547])
	req.SetType(0x0002)
	req.SetServiceMethod(ServiceMethodAuthLogin)
	return tlvs
}

func (h *Handler) checkCaptchaAndGetSessionTickets(
	ctx context.Context,
	req *checkCaptchaAndGetSessionTicketsRequest,
) (*Response, error) {
	return h.getSessionTickets(ctx, req)
}
