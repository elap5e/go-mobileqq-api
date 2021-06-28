package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
)

type AuthGetSessionTicketsWithQRSignatureRequest struct {
	AuthGetSessionTicketsWithPasswordRequest
}

func NewAuthGetSessionTicketsWithQRSignatureRequest(uin uint64, password string) *AuthGetSessionTicketsWithQRSignatureRequest {
	return &AuthGetSessionTicketsWithQRSignatureRequest{
		AuthGetSessionTicketsWithPasswordRequest{
			Username: fmt.Sprintf("%d", uin),

			DstAppID:         defaultClientDstAppID,
			SubDstAppID:      defaultClientOpenAppID,
			AppClientVersion: 0x00000000,
			Uin:              uin,
			I2:               0x0000,
			IPv4Address:      defaultDeviceIPv4Address,
			ServerTime:       0,
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

func (c *Client) AuthGetSessionTicketsWithQRSignature(ctx context.Context, req *AuthGetSessionTicketsWithQRSignatureRequest) (*AuthGetSessionTicketsResponse, error) {
	req.Seq = c.getNextSeq()
	req.TGTGTKey = c.tgtgtKey
	req.T106 = []byte{}
	req.T16A = []byte{}
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
		Type:          0x0009,
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
		Cookie:       c.cookie[:],
		Buffer:       buf,
		ReserveField: c.ksid,
		Simple:       false,
	}, s2c); err != nil {
		return nil, err
	}
	return c.AuthGetSessionTickets(ctx, s2c)
}