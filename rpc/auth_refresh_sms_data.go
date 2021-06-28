package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthRefreshSMSDataRequest struct {
	Seq uint32

	Uin      uint64
	Username string

	Session      []byte
	SMSAppID     uint64
	T174         []byte
	MiscBitmap   uint32
	SubSigMap    uint32
	SubAppIDList []uint64
	SMSExtraData []byte

	lockType uint8
}

func NewAuthRefreshSMSDataRequest(uin uint64) *AuthRefreshSMSDataRequest {
	return &AuthRefreshSMSDataRequest{
		Uin:      uin,
		Username: fmt.Sprintf("%d", uin),

		Session:      nil,
		SMSAppID:     defaultClientSMSAppID,
		T174:         nil,
		MiscBitmap:   clientMiscBitmap,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		SMSExtraData: nil,

		lockType: 0x00,
	}
}

func (req *AuthRefreshSMSDataRequest) GetTLVs(ctx context.Context) (map[uint16]tlv.TLVCodec, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(req.Session)
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0174] = tlv.NewT174(req.T174)
	tlvs[0x017a] = tlv.NewT17A(req.SMSAppID)
	tlvs[0x0197] = tlv.NewTLV(0x0197, 0x0000, bytes.NewBuffer([]byte{req.lockType}))
	tlvs[0x0542] = tlv.NewT542(req.SMSExtraData)
	return tlvs, nil
}

func (c *Client) AuthRefreshSMSData(ctx context.Context, req *AuthRefreshSMSDataRequest) (*AuthGetSessionTicketsResponse, error) {
	req.Seq = c.getNextSeq()
	req.Session = c.session
	req.T174 = c.t174
	tlvs, err := req.GetTLVs(ctx)
	if err != nil {
		return nil, err
	}
	buf, err := oicq.Marshal(ctx, &oicq.Message{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: 0x87,
		RandomKey:     c.randomKey,
		KeyVersion:    c.serverPublicKeyVersion,
		PublicKey:     c.privateKey.Public().Bytes(),
		ShareKey:      c.privateKey.ShareKey(c.serverPublicKey),
		Type:          0x0008,
		TLVs:          tlvs,
	})
	if err != nil {
		return nil, err
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call(ServiceMethodAuthLogin, &ClientToServerMessage{
		Username: req.Username,
		Seq:      req.Seq,
		Buffer:   buf,
		Simple:   false,
	}, s2c); err != nil {
		return nil, err
	}
	return c.AuthGetSessionTickets(ctx, s2c)
}
