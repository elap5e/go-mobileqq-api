package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/rpc/message"
)

type AuthCheckPictureRequest struct {
	*AuthCheckWebSignatureRequest
}

func NewAuthCheckPictureRequest(uin uint64, code, sign []byte) *AuthCheckPictureRequest {
	return &AuthCheckPictureRequest{
		&AuthCheckWebSignatureRequest{
			Uin:      uin,
			Username: fmt.Sprintf("%d", uin),
			CheckWeb: false,

			T104:         nil,
			Code:         code,
			Sign:         sign,
			MiscBitmap:   defaultClientMiscBitmap,
			SubSigMap:    defaultClientSubSigMap,
			SubAppIDList: defaultClientSubAppIDList,
			T547:         nil,
		},
	}
}

func (req *AuthCheckPictureRequest) Encode(ctx context.Context) (*ClientToServerMessage, error) {
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
		Buffer:   buf,
		Simple:   false,
	}, nil
}

func (c *Client) AuthCheckPicture(ctx context.Context, req *AuthCheckPictureRequest) (interface{}, error) {
	req.Seq = c.getNextSeq()
	req.T104 = []byte{}
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
