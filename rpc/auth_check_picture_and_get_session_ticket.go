package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/rpc/message"
)

type AuthCheckPictureAndGetSessionTicketRequest struct {
	*AuthCheckCaptchaAndGetSessionTicketRequest
}

func NewAuthCheckPictureAndGetSessionTicketRequest(uin uint64, code, sign []byte) *AuthCheckPictureAndGetSessionTicketRequest {
	return &AuthCheckPictureAndGetSessionTicketRequest{
		&AuthCheckCaptchaAndGetSessionTicketRequest{
			Uin:      uin,
			Username: fmt.Sprintf("%d", uin),
			CheckWeb: false,

			T104:         nil,
			Code:         code,
			Sign:         sign,
			MiscBitmap:   clientMiscBitmap,
			SubSigMap:    defaultClientSubSigMap,
			SubAppIDList: defaultClientSubAppIDList,
			T547:         nil,
		},
	}
}

func (req *AuthCheckPictureAndGetSessionTicketRequest) Encode(ctx context.Context) (*ClientToServerMessage, error) {
	msg, err := req.EncodeOICQMessage(ctx)
	if err != nil {
		return nil, err
	}
	buf, err := message.MarshalOICQMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &ClientToServerMessage{
		Username: req.Username,
		Seq:      req.Seq,
		AppID:    clientAppID,
		Cookie:   req.Cookie,
		Buffer:   buf,
		Simple:   false,
	}, nil
}

func (c *Client) AuthCheckPictureAndGetSessionTicket(ctx context.Context, req *AuthCheckPictureAndGetSessionTicketRequest) (*AuthGetSessionTicketResponse, error) {
	req.Seq = c.getNextSeq()
	req.Cookie = c.cookie[:]
	req.T104 = c.t104
	c2s, err := req.Encode(ctx)
	if err != nil {
		return nil, err
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call("wtlogin.login", c2s, s2c); err != nil {
		return nil, err
	}
	return c.AuthGetSessionTicket(ctx, s2c)
}
