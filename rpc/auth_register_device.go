package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/rpc/message"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthRegisterDeviceRequest struct {
	Username     string
	Uin          uint64
	Seq          uint32
	SubAppIDList []uint64
	T104         []byte
	T401         [16]byte
}

func NewAuthRegisterDeviceRequest(uin uint64, t104 []byte, t401 [16]byte) *AuthRegisterDeviceRequest {
	return &AuthRegisterDeviceRequest{
		Username:     fmt.Sprintf("%d", uin),
		Uin:          uin,
		SubAppIDList: defaultClientSubAppIDList,
		T104:         t104,
		T401:         t401, // md5.Sum(append(append(defaultDeviceGUID, defaultDeviceDPWD...), t402...))
	}
}

func (req *AuthRegisterDeviceRequest) Marshal(ctx context.Context) ([]byte, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(req.T104)
	tlvs[0x0116] = tlv.NewT116(defaultClientMiscBitmap, defaultClientSubSigMap, req.SubAppIDList)
	tlvs[0x0401] = tlv.NewT401(req.T401)

	return message.MarshalOICQMessage(ctx, &message.OICQMessage{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: 0x07,
		RandomKey:     defaultClientRandomKey,
		PublicKey:     ecdh.PublicKey,
		ShareKey:      ecdh.ShareKey,
		Type:          0x0014,
		TLVs:          tlvs,
	})
}

func (c *Client) AuthRegisterDevice(ctx context.Context, req *AuthRegisterDeviceRequest) (interface{}, error) {
	s2c := new(ServerToClientMessage)
	req.Seq = c.getNextSeq()
	buf, err := req.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.Call("wtlogin.login", &ClientToServerMessage{
		Username: req.Username,
		Seq:      req.Seq,
		Buffer:   buf,
		Simple:   true,
	}, s2c); err != nil {
		return nil, err
	}
	return s2c, nil
}
