package auth

import (
	"context"
)

type getSessionTicketsWithQRSignatureRequest struct {
	getSessionTicketsWithPasswordRequest
}

func newGetSessionTicketsWithQRSignatureRequest(
	username string,
	password string,
) *getSessionTicketsWithQRSignatureRequest {
	req := &getSessionTicketsWithQRSignatureRequest{
		getSessionTicketsWithPasswordRequest{
			DstAppID:         defaultClientDstAppID,
			SubDstAppID:      defaultClientOpenAppID,
			AppClientVersion: 0x00000000,
			_Uin:             0x00000000,
			I2:               0x0000,
			_IPv4Address:     nil,        // nil
			ServerTime:       0,          // nil
			PasswordMD5:      [16]byte{}, // nil
			_UserA1Key:       [16]byte{}, // nil
			LoginType:        0x00000000,
			UserA1:           nil,
			T16A:             nil,
			_MiscBitmap:      0x00000000,
			SubSigMap:        defaultClientSubSigMap,
			SubAppIDList:     defaultClientSubAppIDList,
			MainSigMap:       defaultClientMainSigMap & 0xfdfffffe,
			SrcAppID:         defaultClientOpenAppID,
			I7:               0x0000,
			I8:               0x00,
			I9:               0x0000,
			I10:              0x01,
			_KSID:            nil,
			_AuthSession:     nil,
			_PackageName:     []byte{},
			Domains:          defaultClientDomains,
		},
	}
	req.SetUsername(username)
	return req
}

func (h *Handler) AuthGetSessionTicketsWithQRSignature(
	ctx context.Context,
	req *getSessionTicketsWithQRSignatureRequest,
) (*Response, error) {
	return h.getSessionTickets(ctx, req)
}
