package rpc

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthGetSessionTicketsWithoutPasswordRequest struct {
	Seq  uint32
	T172 []byte

	Username string

	Uin              uint64
	DstAppID         uint64
	SrcAppID         uint64
	AppClientVersion uint32
	MainSigMap       uint32
	UserA2           []byte
	MiscBitmap       uint32
	SubSigMap        uint32
	SubAppIDList     []uint64
	KSID             []byte
	UserD2           []byte
	Domains          []string
}

func NewAuthGetSessionTicketsWithoutPasswordRequest(uin uint64) *AuthGetSessionTicketsWithoutPasswordRequest {
	return &AuthGetSessionTicketsWithoutPasswordRequest{
		Username: fmt.Sprintf("%d", uin),

		Uin:              uin,
		DstAppID:         defaultClientDstAppID,
		SrcAppID:         defaultClientOpenAppID,
		AppClientVersion: 0x00000000,
		MainSigMap:       defaultClientMainSigMap & 0xfdfffffe,
		UserA2:           nil,
		MiscBitmap:       clientMiscBitmap,
		SubSigMap:        defaultClientSubSigMap,
		SubAppIDList:     defaultClientSubAppIDList,
		KSID:             GetClientCodecKSID(),
		UserD2:           nil,
		Domains:          defaultClientDomains,
	}
}

func (req *AuthGetSessionTicketsWithoutPasswordRequest) GetTLVs(ctx context.Context) (map[uint16]tlv.TLVCodec, error) {
	key := SelectClientCodecKey(req.Username)
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0100] = tlv.NewT100(req.DstAppID, req.SrcAppID, req.AppClientVersion, req.MainSigMap)
	tlvs[0x010a] = tlv.NewT10A(req.UserA2)
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0108] = tlv.NewT108(req.KSID)
	tlvs[0x0144] = tlv.NewT144(md5.Sum(key.D2Key[:]),
		tlv.NewT109(md5.Sum(defaultDeviceOSBuildID)),
		tlv.NewT52D(ctx),
		tlv.NewT124([]byte(defaultDeviceOSType), []byte(defaultDeviceOSVersion), defaultDeviceNetworkTypeID, defaultDeviceSIMOPName, nil, defaultDeviceAPNName),
		tlv.NewT128(deviceIsGUIDFileNil, deviceIsGUIDGenSucc, deviceIsGUIDChanged, deviceGUIDFlag, []byte(defaultDeviceOSBuildModel), deviceGUID[:], defaultDeviceOSBuildBrand),
		tlv.NewT16E([]byte(defaultDeviceOSBuildModel)),
	)
	tlvs[0x0143] = tlv.NewT143(req.UserD2)
	tlvs[0x0142] = tlv.NewT142(clientPackageName)
	tlvs[0x0154] = tlv.NewT154(req.Seq)
	tlvs[0x0018] = tlv.NewT18(req.DstAppID, req.MainSigMap, req.Uin, 0x0000)
	tlvs[0x0141] = tlv.NewT141(defaultDeviceSIMOPName, defaultDeviceNetworkTypeID, defaultDeviceAPNName)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	if len(req.Domains) > 0 {
		tlvs[0x0511] = tlv.NewT511(req.Domains)
	}
	tlvs[0x0147] = tlv.NewT147(req.DstAppID, clientVersionName, clientSignatureMD5)
	if len(req.T172) != 0 {
		tlvs[0x0172] = tlv.NewT172(req.T172)
	}
	tlvs[0x0177] = tlv.NewT177(clientBuildTime, clientSDKVersion)
	tlvs[0x0187] = tlv.NewT187(md5.Sum(defaultDeviceMACAddress))
	tlvs[0x0188] = tlv.NewT188(md5.Sum(defaultDeviceOSBuildID))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(defaultDeviceIMSI)))
	// DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(md5.Sum(defaultDeviceBSSIDAddress), defaultDeviceSSIDAddress)
	// DISABLED: tgt
	// tlvs[0x0544] = tlv.NewT544(req.Uin, deviceGUID, clientSDKVersion, 0x0009)
	return tlvs, nil
}

func (c *Client) AuthGetSessionTicketsWithoutPassword(ctx context.Context, req *AuthGetSessionTicketsWithoutPasswordRequest) (*AuthGetSessionTicketsResponse, error) {
	req.Seq = c.getNextSeq()
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
		Type:          0x000b,
		TLVs:          tlvs,
	})
	if err != nil {
		return nil, err
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call(ServiceMethodAuthExchangeAccount, &ClientToServerMessage{
		Username: req.Username,
		Seq:      req.Seq,
		Buffer:   buf,
		Simple:   false,
	}, s2c); err != nil {
		return nil, err
	}
	return c.AuthGetSessionTickets(ctx, s2c)
}
