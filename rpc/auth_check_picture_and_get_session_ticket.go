package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
)

type AuthCheckPictureAndGetSessionTicketRequest struct {
	AuthCheckCaptchaAndGetSessionTicketRequest
}

func NewAuthCheckPictureAndGetSessionTicketRequest(uin uint64, code, sign []byte) *AuthCheckPictureAndGetSessionTicketRequest {
	return &AuthCheckPictureAndGetSessionTicketRequest{
		AuthCheckCaptchaAndGetSessionTicketRequest{
			Uin:      uin,
			Username: fmt.Sprintf("%d", uin),

			T104:         nil,
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

func (c *Client) AuthCheckPictureAndGetSessionTicket(ctx context.Context, req *AuthCheckPictureAndGetSessionTicketRequest) (*AuthGetSessionTicketResponse, error) {
	req.Seq = c.getNextSeq()
	req.Cookie = c.cookie[:]
	req.T104 = c.t104
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
