package client

import (
	"context"
)

type AuthCheckPictureAndGetSessionTicketsRequest struct {
	AuthCheckCaptchaAndGetSessionTicketsRequest
}

func NewAuthCheckPictureAndGetSessionTicketsRequest(
	username string,
	code, sign []byte,
) *AuthCheckPictureAndGetSessionTicketsRequest {
	req := &AuthCheckPictureAndGetSessionTicketsRequest{
		AuthCheckCaptchaAndGetSessionTicketsRequest{
			_AuthSession: nil,
			Code:         code,
			Sign:         sign,
			_MiscBitmap:  0x00000000,
			SubSigMap:    defaultClientSubSigMap,
			SubAppIDList: defaultClientSubAppIDList,
			_ExtraData:   nil,

			isCaptcha: false,
		},
	}
	req.SetUsername(username)
	return req
}

func (c *Client) AuthCheckPictureAndGetSessionTickets(
	ctx context.Context,
	req *AuthCheckPictureAndGetSessionTicketsRequest,
) (*AuthGetSessionTicketsResponse, error) {
	return c.AuthGetSessionTickets(ctx, req)
}
