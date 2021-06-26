package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
)

type AuthGetSessionTicketWithQRSignatureRequest struct {
	*AuthGetSessionTicketWithPasswordRequest
}

func NewAuthGetSessionTicketWithQRSignatureRequest(uin uint64, password string) *AuthGetSessionTicketWithQRSignatureRequest {
	return &AuthGetSessionTicketWithQRSignatureRequest{
		&AuthGetSessionTicketWithPasswordRequest{
			Username: fmt.Sprintf("%d", uin),

			DstAppID:         defaultClientDstAppID,
			SubDstAppID:      defaultClientOpenAppID,
			AppClientVersion: 0x00000000,
			Uin:              uin,
			I2:               0x0000,
			IPv4Address:      defaultDeviceIPv4Address,
			CurrentTime:      0,
			PasswordMD5:      [16]byte{},
			TGTGTKey:         [16]byte{},
			LoginType:        0x00000000,
			T106:             nil,
			T16A:             nil,
			MiscBitmap:       clientMiscBitmap,
			SubSigMap:        defaultClientSubSigMap,
			SubAppIDList:     defaultClientSubAppIDList,
			MainSigMap:       defaultClientMainSigMap & 0xfdfffffe,
			SrcAppID:         defaultClientOpenAppID,
			I7:               0x0000,
			I8:               0x00,
			I9:               0x0000,
			I10:              0x01,
			KSID:             GetClientCodecKSID(),
			T104:             nil,
			PackageName:      clientPackageName,
			Domains:          defaultClientDomains,
		},
	}
}

func (req *AuthGetSessionTicketWithQRSignatureRequest) Encode(ctx context.Context) (*ClientToServerMessage, error) {
	msg, err := req.EncodeOICQMessage(ctx)
	if err != nil {
		return nil, err
	}
	buf, err := oicq.Marshal(ctx, msg)
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

func (c *Client) AuthGetSessionTicketWithQRSignature(ctx context.Context, req *AuthGetSessionTicketWithQRSignatureRequest) (*AuthGetSessionTicketResponse, error) {
	req.Seq = c.getNextSeq()
	req.TGTGTKey = c.tgtgtKey
	req.T106 = []byte{}
	req.T16A = []byte{}
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
