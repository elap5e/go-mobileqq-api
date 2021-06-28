package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
)

type AuthCheckPictureAndGetSessionTicketsRequest struct {
	AuthCheckCaptchaAndGetSessionTicketsRequest
}

func NewAuthCheckPictureAndGetSessionTicketsRequest(uin uint64, code, sign []byte) *AuthCheckPictureAndGetSessionTicketsRequest {
	return &AuthCheckPictureAndGetSessionTicketsRequest{
		AuthCheckCaptchaAndGetSessionTicketsRequest{
			Uin:      uin,
			Username: fmt.Sprintf("%d", uin),

			Session:      nil,
			Code:         code,
			Sign:         sign,
			MiscBitmap:   clientMiscBitmap,
			SubSigMap:    defaultClientSubSigMap,
			SubAppIDList: defaultClientSubAppIDList,
			T547:         nil,

			checkWeb: false,
		},
	}
}

func (c *Client) AuthCheckPictureAndGetSessionTickets(ctx context.Context, req *AuthCheckPictureAndGetSessionTicketsRequest) (*AuthGetSessionTicketsResponse, error) {
	req.Seq = c.getNextSeq()
	req.Cookie = c.cookie[:]
	req.Session = c.session
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
