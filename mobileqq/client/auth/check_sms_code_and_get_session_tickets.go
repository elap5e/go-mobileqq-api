package auth

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type checkSMSCodeAndGetSessionTicketsRequest struct {
	request

	_AuthSession []byte // h.GetUserSignature(req.Username).Session.Auth
	Code         []byte
	_T174        []byte // h.t174
	_MiscBitmap  uint32 // h.opt.Client.MiscBitmap
	SubSigMap    uint32
	SubAppIDList []uint64
	_HashedGUID  [16]byte // h.hashedGUID
	_ExtraData   []byte   // h.extraData[0x0542]

	lockType uint8
}

func newCheckSMSCodeAndGetSessionTicketsRequest(
	username string,
	code []byte,
) *checkSMSCodeAndGetSessionTicketsRequest {
	req := &checkSMSCodeAndGetSessionTicketsRequest{
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

func (req *checkSMSCodeAndGetSessionTicketsRequest) MustGetTLVs(
	ctx context.Context,
) map[uint16]tlv.TLVCodec {
	h := forHandler(ctx)
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(
		h.client.GetUserSignature(req.GetUsername()).Session.Auth,
	)
	tlvs[0x0116] = tlv.NewT116(
		h.opt.Client.MiscBitmap,
		req.SubSigMap,
		req.SubAppIDList,
	)
	tlvs[0x0174] = tlv.NewT174(h.t174)
	tlvs[0x017c] = tlv.NewT17C(req.Code)
	tlvs[0x0401] = tlv.NewT401(h.hashedGUID)
	tlvs[0x0197] = tlv.NewTLV(
		0x0198,
		0x0000,
		bytes.NewBuffer([]byte{req.lockType}),
	)
	tlvs[0x0542] = tlv.NewT542(h.extraData[0x0542])
	req.SetType(0x0007)
	req.SetServiceMethod(ServiceMethodAuthLogin)
	return tlvs
}

func (h *Handler) checkSMSCodeAndGetSessionTickets(
	ctx context.Context,
	req *checkSMSCodeAndGetSessionTicketsRequest,
) (*Response, error) {
	return h.getSessionTickets(ctx, req)
}
