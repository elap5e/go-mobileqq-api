package rpc

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthRefreshSMSDataRequest struct {
	authGetSessionTicketsRequest

	AuthSession  []byte // c.GetUserSignature(req.Username).Session.Auth
	SMSAppID     uint64
	T174         []byte // c.t174
	MiscBitmap   uint32 // c.cfg.Client.MiscBitmap
	SubSigMap    uint32
	SubAppIDList []uint64
	SMSExtraData []byte // c.extraData[0x0542]

	lockType uint8
}

func NewAuthRefreshSMSDataRequest(
	username string,
) *AuthRefreshSMSDataRequest {
	req := &AuthRefreshSMSDataRequest{
		AuthSession:  nil,
		SMSAppID:     defaultClientSMSAppID,
		T174:         nil,
		MiscBitmap:   0x00000000,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		SMSExtraData: nil,

		lockType: 0x00,
	}
	req.SetUsername(username)
	return req
}

func (req *AuthRefreshSMSDataRequest) GetTLVs(
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
	tlvs[0x017a] = tlv.NewT17A(req.SMSAppID)
	tlvs[0x0197] = tlv.NewTLV(
		0x0197,
		0x0000,
		bytes.NewBuffer([]byte{req.lockType}),
	)
	tlvs[0x0542] = tlv.NewT542(c.extraData[0x0542])
	req.SetType(0x0008)
	return tlvs, nil
}

func (c *Client) AuthRefreshSMSData(
	ctx context.Context,
	req *AuthRefreshSMSDataRequest,
) (*AuthGetSessionTicketsResponse, error) {
	return c.AuthGetSessionTickets(ctx, req)
}
