package rpc

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/rpc/message"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthGetSessionTicketWithoutPasswordRequest struct {
	Username     string
	Uin          uint64
	DstAppID     uint64
	MainSigMap   uint32
	SubDstAppID  uint64
	SubAppIDList []uint64
	Domains      []string
	AuthData     *ClientAuthData
}

func NewAuthGetSessionTicketWithoutPasswordRequest(uin uint64, auth *ClientAuthData) *AuthGetSessionTicketWithoutPasswordRequest {
	return &AuthGetSessionTicketWithoutPasswordRequest{
		Username:     fmt.Sprintf("%d", uin),
		Uin:          uin,
		DstAppID:     defaultClientDstAppID,
		MainSigMap:   defaultClientMainSigMap,
		SubDstAppID:  defaultClientSubDstAppID,
		SubAppIDList: defaultClientSubAppIDList,
		Domains:      defaultClientDomains,
		AuthData:     auth,
	}
}

func (req *AuthGetSessionTicketWithoutPasswordRequest) Marshal(ctx context.Context) ([]byte, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0100] = tlv.NewT100(req.DstAppID, req.SubDstAppID, 0x00000000, req.MainSigMap&0xfdfffffe)
	tlvs[0x010a] = tlv.NewT10A(req.AuthData.A2)
	tlvs[0x0116] = tlv.NewT116(defaultClientMiscBitmap, defaultClientSubSigMap, req.SubAppIDList)
	tlvs[0x0108] = tlv.NewT108(defaultDeviceKSID)
	tlvs[0x0144] = tlv.NewT144(md5.Sum(req.AuthData.Key[:]),
		tlv.NewT109(md5.Sum(defaultDeviceOSID)),
		tlv.NewT52D(ctx),
		tlv.NewT124(defaultDeviceOSType, defaultDeviceOSBuildVersionRelease, defaultDeviceNetworkTypeID, defaultDeviceSIMOPName, nil, defaultDeviceAPNName),
		tlv.NewT128(false, true, false, (1<<24&0xFF000000)|(0<<8&0xFF00), defaultDeviceOSBuildModel, defaultDeviceGUID, defaultDeviceOSBuildBrand),
		tlv.NewT16E(defaultDeviceOSBuildModel),
	)
	tlvs[0x0143] = tlv.NewT143(req.AuthData.D2)
	tlvs[0x0142] = tlv.NewT142(defaultAPKID)
	tlvs[0x0154] = tlv.NewT154(0x0000)
	tlvs[0x0018] = tlv.NewT18(req.DstAppID, req.MainSigMap&0xfdfffffe, req.Uin, 0x0000)
	tlvs[0x0141] = tlv.NewT141(defaultDeviceSIMOPName, defaultDeviceNetworkTypeID, defaultDeviceAPNName)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	if len(req.Domains) > 0 {
		tlvs[0x0511] = tlv.NewT511(req.Domains)
	}
	tlvs[0x0147] = tlv.NewT147(req.DstAppID, defaultAPKVersionName, defaultAPKSignatureMD5)
	// tlvs[0x0172] = tlv.NewT172([]byte{})
	tlvs[0x0177] = tlv.NewT177(defaultClientBuildTime, defaultClientSDKVersion)
	tlvs[0x0187] = tlv.NewT187(md5.Sum(defaultDeviceMACAddress))
	tlvs[0x0188] = tlv.NewT188(md5.Sum(defaultDeviceOSID))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(defaultDeviceIMSI)))
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(md5.Sum(defaultDeviceBSSIDAddress), defaultDeviceSSIDAddress)
	// tlvs[0x0544] = tlv.NewT544(req.Username, "810_a", nil)

	return message.MarshalOICQMessage(ctx, &message.OICQMessage{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: 0x07,
		RandomKey:     defaultClientRandomKey,
		PublicKey:     ecdh.PublicKey,
		ShareKey:      ecdh.ShareKey,
		Type:          0x000b,
		TLVs:          tlvs,
	})
}

func (req *AuthGetSessionTicketWithoutPasswordRequest) Unmarshal(ctx context.Context, msg *message.OICQMessage) error {
	return nil
}

func (c *Client) AuthGetSessionTicketWithoutPassword(ctx context.Context, req *AuthGetSessionTicketWithoutPasswordRequest) (interface{}, error) {
	s2c := new(ServerToClientMessage)
	buf, err := req.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.Call("wtlogin.exchange_emp", &ClientToServerMessage{
		Username: req.Username,
		Seq:      c.getNextSeq(),
		Buffer:   buf,
		Simple:   true,
	}, s2c); err != nil {
		return nil, err
	}
	return s2c, nil
}
