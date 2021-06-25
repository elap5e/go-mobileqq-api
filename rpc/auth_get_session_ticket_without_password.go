package rpc

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/rpc/message"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthGetSessionTicketWithoutPasswordRequest struct {
	Username string
	Seq      uint32

	Uin              uint64
	DstAppID         uint64
	SrcAppID         uint64
	AppClientVersion uint32
	MainSigMap       uint32
	T10A             []byte
	MiscBitmap       uint32
	SubSigMap        uint32
	SubAppIDList     []uint64
	KSID             []byte
	T143             []byte
	Domains          []string
}

func NewAuthGetSessionTicketWithoutPasswordRequest(uin uint64) *AuthGetSessionTicketWithoutPasswordRequest {
	return &AuthGetSessionTicketWithoutPasswordRequest{
		Username: fmt.Sprintf("%d", uin),

		Uin:              uin,
		DstAppID:         defaultClientDstAppID,
		SrcAppID:         defaultClientOpenAppID,
		AppClientVersion: 0x00000000,
		MainSigMap:       defaultClientMainSigMap & 0xfdfffffe,
		T10A:             nil,
		MiscBitmap:       clientMiscBitmap,
		SubSigMap:        defaultClientSubSigMap,
		SubAppIDList:     defaultClientSubAppIDList,
		KSID:             GetClientCodecKSID(),
		T143:             nil,
		Domains:          defaultClientDomains,
	}
}

func (req *AuthGetSessionTicketWithoutPasswordRequest) Marshal(ctx context.Context) ([]byte, error) {
	key := SelectClientCodecKey(req.Username)
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0100] = tlv.NewT100(req.DstAppID, req.SrcAppID, req.AppClientVersion, req.MainSigMap)
	tlvs[0x010a] = tlv.NewT10A(key.A2)
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0108] = tlv.NewT108(req.KSID)
	tlvs[0x0144] = tlv.NewT144(md5.Sum(key.Key[:]),
		tlv.NewT109(md5.Sum(defaultDeviceOSBuildID)),
		tlv.NewT52D(ctx),
		tlv.NewT124(defaultDeviceOSType, defaultDeviceOSVersion, defaultDeviceNetworkTypeID, defaultDeviceSIMOPName, nil, defaultDeviceAPNName),
		tlv.NewT128(deviceIsGUIDFileNil, deviceIsGUIDGenSucc, deviceIsGUIDChanged, deviceGUIDFlag, defaultDeviceOSBuildModel, deviceGUID[:], defaultDeviceOSBuildBrand),
		tlv.NewT16E(defaultDeviceOSBuildModel),
	)
	tlvs[0x0143] = tlv.NewT143(key.D2)
	tlvs[0x0142] = tlv.NewT142(clientPackageName)
	tlvs[0x0154] = tlv.NewT154(req.Seq)
	tlvs[0x0018] = tlv.NewT18(req.DstAppID, req.MainSigMap, req.Uin, 0x0000)
	tlvs[0x0141] = tlv.NewT141(defaultDeviceSIMOPName, defaultDeviceNetworkTypeID, defaultDeviceAPNName)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	if len(req.Domains) > 0 {
		tlvs[0x0511] = tlv.NewT511(req.Domains)
	}
	tlvs[0x0147] = tlv.NewT147(req.DstAppID, clientVersionName, clientSignatureMD5)
	// tlvs[0x0172] = tlv.NewT172([]byte{})
	tlvs[0x0177] = tlv.NewT177(clientBuildTime, clientSDKVersion)
	tlvs[0x0187] = tlv.NewT187(md5.Sum(defaultDeviceMACAddress))
	tlvs[0x0188] = tlv.NewT188(md5.Sum(defaultDeviceOSBuildID))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(defaultDeviceIMSI)))
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(md5.Sum(defaultDeviceBSSIDAddress), defaultDeviceSSIDAddress)
	// tlvs[0x0544] = tlv.NewT544(req.Username, "810_a", nil)

	return message.MarshalOICQMessage(ctx, &message.OICQMessage{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: 0x07,
		RandomKey:     clientRandomKey,
		KeyVersion:    ecdh.KeyVersion,
		PublicKey:     ecdh.PublicKey,
		ShareKey:      ecdh.ShareKey,
		Type:          0x000b,
		TLVs:          tlvs,
	})
}

func (c *Client) AuthGetSessionTicketWithoutPassword(ctx context.Context, req *AuthGetSessionTicketWithoutPasswordRequest) (interface{}, error) {
	s2c := new(ServerToClientMessage)
	req.Seq = c.getNextSeq()
	buf, err := req.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.Call("wtlogin.exchange_emp", &ClientToServerMessage{
		Username: req.Username,
		Seq:      req.Seq,
		AppID:    clientAppID,
		Buffer:   buf,
		Simple:   false,
	}, s2c); err != nil {
		return nil, err
	}
	return c.AuthGetSessionTicket(ctx, s2c)
}
