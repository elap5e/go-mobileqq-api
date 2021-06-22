package rpc

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/rpc/message"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthGetSessionTicketWithPasswordRequest struct {
	Username     string
	Uin          uint64
	Seq          uint32
	DstAppID     uint64
	MainSigMap   uint32
	SubDstAppID  uint64
	PasswordMD5  [16]byte
	SubAppIDList []uint64
	Domains      []string
	T104         []byte
}

func NewAuthGetSessionTicketWithPasswordRequest(uin uint64, password string) *AuthGetSessionTicketWithPasswordRequest {
	return &AuthGetSessionTicketWithPasswordRequest{
		Username:     fmt.Sprintf("%d", uin),
		Uin:          uin,
		DstAppID:     defaultClientDstAppID,
		MainSigMap:   defaultClientMainSigMap,
		SubDstAppID:  defaultClientSubDstAppID,
		PasswordMD5:  md5.Sum([]byte(password)),
		SubAppIDList: defaultClientSubAppIDList,
		Domains:      defaultClientDomains,
	}
}

func (req *AuthGetSessionTicketWithPasswordRequest) Marshal(ctx context.Context) ([]byte, error) {
	var tgtgtKey [16]byte
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0018] = tlv.NewT18(req.DstAppID, req.MainSigMap&0xfdfffffe, req.Uin, 0x0000)
	tlvs[0x0001] = tlv.NewT1(req.Uin, defaultDeviceIPv4Address)
	tlvs[0x0106] = tlv.NewT106(req.DstAppID, req.SubDstAppID, 0x00000000, req.Uin, nil, defaultDeviceIPv4Address, true, req.PasswordMD5, 0, []byte(req.Username), tgtgtKey, true, defaultDeviceGUID, 1)
	tlvs[0x0116] = tlv.NewT116(defaultClientMiscBitmap, defaultClientSubSigMap, req.SubAppIDList)
	tlvs[0x0100] = tlv.NewT100(req.DstAppID, req.SubDstAppID, 0x00000000, req.MainSigMap&0xfdfffffe)
	tlvs[0x0107] = tlv.NewT107(0x0000, 0x00, 0x0000, 0x01)
	if len(defaultDeviceKSID) != 0 {
		tlvs[0x0108] = tlv.NewT108(defaultDeviceKSID)
	}
	if len(req.T104) != 0 {
		tlvs[0x0104] = tlv.NewT104(req.T104)
	}
	tlvs[0x0142] = tlv.NewT142(defaultAPKID)
	// tlvs[0x0112] = tlv.NewT112(nil)
	tlvs[0x0144] = tlv.NewT144(tgtgtKey,
		tlv.NewT109(md5.Sum(defaultDeviceOSID)),
		tlv.NewT52D(ctx),
		tlv.NewT124(defaultDeviceOSType, defaultDeviceOSBuildVersionRelease, defaultDeviceNetworkTypeID, defaultDeviceSIMOPName, nil, defaultDeviceAPNName),
		tlv.NewT128(false, true, false, (1<<24&0xFF000000)|(0<<8&0xFF00), defaultDeviceOSBuildModel, defaultDeviceGUID, defaultDeviceOSBuildBrand),
		tlv.NewT16E(defaultDeviceOSBuildModel),
	)
	tlvs[0x0145] = tlv.NewT145(defaultDeviceGUID)
	tlvs[0x0147] = tlv.NewT147(req.DstAppID, defaultAPKVersionName, defaultAPKSignatureMD5)
	tlvs[0x0166] = tlv.NewT166(0x01)
	// tlvs[0x016a] = tlv.NewT16A(nil)
	tlvs[0x0154] = tlv.NewT154(req.Seq)
	tlvs[0x0141] = tlv.NewT141(defaultDeviceSIMOPName, defaultDeviceNetworkTypeID, defaultDeviceAPNName)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	if len(req.Domains) > 0 {
		tlvs[0x0511] = tlv.NewT511(req.Domains)
	}
	// tlvs[0x0172] = tlv.NewT172([]byte{})
	tlvs[0x0185] = tlv.NewT185(0x01)
	// tlvs[0x0400] = tlv.NewT400([]byte{})
	tlvs[0x0187] = tlv.NewT187(md5.Sum(defaultDeviceMACAddress))
	tlvs[0x0188] = tlv.NewT188(md5.Sum(defaultDeviceOSID))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(defaultDeviceIMSI)))
	tlvs[0x0191] = tlv.NewT191(0x82)
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(md5.Sum(defaultDeviceBSSIDAddress), defaultDeviceSSIDAddress)
	tlvs[0x0177] = tlv.NewT177(defaultClientBuildTime, defaultClientSDKVersion)
	tlvs[0x0516] = tlv.NewTLV(0x0516, 0x0004, bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00}))
	tlvs[0x0521] = tlv.NewTLV(0x0521, 0x0006, bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
	// buf := bytes.NewBuffer([]byte{})
	// tlv.NewTLV(0x0536, 0x0002, bytes.NewBuffer([]byte{0x01, 0x00})).Encode(buf)
	// tlvs[0x0525] = tlv.NewTLV(0x0525, 0x0000, buf)
	tlvs[0x0525] = tlv.NewTLV(0x0525, 0x0006, bytes.NewBuffer([]byte{0x05, 0x25, 0x00, 0x02, 0x01, 0x00}))
	// tlvs[0x0318] = tlv.NewT318(nil)
	// tlvs[0x0544] = tlv.NewT544(req.Username, "810_a", nil)
	tlvs[0x0545] = tlv.NewT545(md5.Sum([]byte("qimei")))
	tlvs[0x0548] = tlv.NewT548([]byte("nativeGetTestData"))

	return message.MarshalOICQMessage(ctx, &message.OICQMessage{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: 0x07,
		RandomKey:     defaultClientRandomKey,
		PublicKey:     ecdh.PublicKey,
		ShareKey:      ecdh.ShareKey,
		Type:          0x0009,
		TLVs:          tlvs,
	})
}

func (c *Client) AuthGetSessionTicketWithPassword(ctx context.Context, req *AuthGetSessionTicketWithPasswordRequest) (interface{}, error) {
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
	return c.AuthGetSessionTicket(ctx, s2c)
}
