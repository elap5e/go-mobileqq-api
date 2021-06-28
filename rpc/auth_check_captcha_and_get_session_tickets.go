package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthCheckCaptchaAndGetSessionTicketsRequest struct {
	Seq    uint32
	Cookie []byte

	Uin      uint64
	Username string

	Session      []byte
	Code         []byte
	Sign         []byte
	MiscBitmap   uint32
	SubSigMap    uint32
	SubAppIDList []uint64
	T547         []byte

	checkWeb bool
}

func NewAuthCheckCaptchaAndGetSessionTicketsRequest(uin uint64, code []byte) *AuthCheckCaptchaAndGetSessionTicketsRequest {
	return &AuthCheckCaptchaAndGetSessionTicketsRequest{
		Uin:      uin,
		Username: fmt.Sprintf("%d", uin),

		Session:      nil,
		Code:         code,
		Sign:         nil,
		MiscBitmap:   clientMiscBitmap,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		T547:         nil,

		checkWeb: true,
	}
}

func (req *AuthCheckCaptchaAndGetSessionTicketsRequest) GetTLVs(ctx context.Context) (map[uint16]tlv.TLVCodec, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	if req.checkWeb {
		tlvs[0x0193] = tlv.NewT193(req.Code)
	} else {
		tlvs[0x0002] = tlv.NewT2(req.Code, req.Sign)
	}
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(req.Session)
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0547] = tlv.NewT547(req.T547)
	return tlvs, nil
}

func (c *Client) AuthCheckCaptchaAndGetSessionTickets(ctx context.Context, req *AuthCheckCaptchaAndGetSessionTicketsRequest) (*AuthGetSessionTicketsResponse, error) {
	req.Seq = c.getNextSeq()
	req.Cookie = c.cookie[:]
	req.Session = c.session
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
		Username: req.Username,
		Seq:      req.Seq,
		Buffer:   buf,
		Simple:   false,
	}, s2c); err != nil {
		return nil, err
	}
	return c.AuthGetSessionTickets(ctx, s2c)
}
