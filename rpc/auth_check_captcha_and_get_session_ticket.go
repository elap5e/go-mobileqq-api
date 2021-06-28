package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthCheckCaptchaAndGetSessionTicketRequest struct {
	Seq    uint32
	Cookie []byte

	Uin      uint64
	Username string

	T104         []byte
	Code         []byte
	Sign         []byte
	MiscBitmap   uint32
	SubSigMap    uint32
	SubAppIDList []uint64
	T547         []byte

	checkWeb bool
}

func NewAuthCheckCaptchaAndGetSessionTicketRequest(uin uint64, code []byte) *AuthCheckCaptchaAndGetSessionTicketRequest {
	return &AuthCheckCaptchaAndGetSessionTicketRequest{
		Uin:      uin,
		Username: fmt.Sprintf("%d", uin),

		T104:         nil,
		Code:         code,
		Sign:         nil,
		MiscBitmap:   clientMiscBitmap,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		T547:         nil,

		checkWeb: true,
	}
}

func (req *AuthCheckCaptchaAndGetSessionTicketRequest) GetTLVs(ctx context.Context) (map[uint16]tlv.TLVCodec, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	if req.checkWeb {
		tlvs[0x0193] = tlv.NewT193(req.Code)
	} else {
		tlvs[0x0002] = tlv.NewT2(req.Code, req.Sign)
	}
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(req.T104)
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0547] = tlv.NewT547(req.T547)
	return tlvs, nil
}

func (c *Client) AuthCheckCaptchaAndGetSessionTicket(ctx context.Context, req *AuthCheckCaptchaAndGetSessionTicketRequest) (*AuthGetSessionTicketResponse, error) {
	req.Seq = c.getNextSeq()
	req.Cookie = c.cookie[:]
	req.T104 = c.t104
	req.T547 = c.t547
	tlvs, err := req.GetTLVs(ctx)
	if err != nil {
		return nil, err
	}
	buf, err := oicq.Marshal(ctx, &oicq.Message{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: oicq.EncryptMethodECDH,
		RandomKey:     c.randomKey,
		KeyVersion:    c.serverPublicKeyVersion,
		PublicKey:     c.privateKey.Public().Bytes(),
		ShareKey:      c.privateKey.ShareKey(c.serverPublicKey),
		Type:          0x0002,
		TLVs:          tlvs,
	})
	if err != nil {
		return nil, err
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call(ServiceMethodAuthLogin, &ClientToServerMessage{
		Username:     req.Username,
		Seq:          req.Seq,
		AppID:        clientAppID,
		Cookie:       req.Cookie,
		Buffer:       buf,
		ReserveField: c.ksid,
		Simple:       false,
	}, s2c); err != nil {
		return nil, err
	}
	return c.AuthGetSessionTicket(ctx, s2c)
}
