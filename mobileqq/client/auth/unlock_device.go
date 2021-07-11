package auth

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/tlv"
)

type unlockDeviceRequest struct {
	request

	_AuthSession []byte // h.GetUserSignature(req.Username).Session.Auth
	_MiscBitmap  uint32 // h.opt.Client.MiscBitmap
	SubSigMap    uint32
	SubAppIDList []uint64
	_HashedGUID  [16]byte // h.hashedGUID
}

func newUnlockDeviceRequest(username string) *unlockDeviceRequest {
	req := &unlockDeviceRequest{
		_AuthSession: nil,
		_MiscBitmap:  0x00000000,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		_HashedGUID:  [16]byte{},
	}
	req.SetUsername(username)
	return req
}

func (req *unlockDeviceRequest) MustGetTLVs(
	ctx context.Context,
) map[uint16]tlv.TLVCodec {
	h := ForHandler(ctx)
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(
		h.GetUserSignature(req.GetUsername()).Session.Auth,
	)
	tlvs[0x0116] = tlv.NewT116(
		h.opt.Client.MiscBitmap,
		req.SubSigMap,
		req.SubAppIDList,
	)
	tlvs[0x0401] = tlv.NewT401(h.hashedGUID)
	req.SetType(0x0014)
	req.SetServiceMethod(ServiceMethodAuthLogin)
	return tlvs
}

func (h *Handler) unlockDevice(
	ctx context.Context,
	req *unlockDeviceRequest,
) (*Response, error) {
	return h.getSessionTickets(ctx, req)
}
