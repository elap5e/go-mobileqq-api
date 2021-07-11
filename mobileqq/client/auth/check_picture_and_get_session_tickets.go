package auth

import (
	"context"
)

type checkPictureAndGetSessionTicketsRequest struct {
	checkCaptchaAndGetSessionTicketsRequest
}

func newCheckPictureAndGetSessionTicketsRequest(
	username string,
	code, sign []byte,
) *checkPictureAndGetSessionTicketsRequest {
	req := &checkPictureAndGetSessionTicketsRequest{
		checkCaptchaAndGetSessionTicketsRequest{
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

func (h *Handler) checkPictureAndGetSessionTickets(
	ctx context.Context,
	req *checkPictureAndGetSessionTicketsRequest,
) (*Response, error) {
	return h.getSessionTickets(ctx, req)
}
