package rpc

import (
	"context"
	"crypto/md5"
	"fmt"
	"net"
	"time"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/rpc/message"
	"github.com/elap5e/go-mobileqq-api/tlv"
	"github.com/elap5e/go-mobileqq-api/util"
)

type AuthGetSessionTicketWithPasswordRequest struct {
	Seq       uint32
	Username  string
	ImageType uint8
	T172      []byte
	T9        []byte
	TGTQR     []byte

	DstAppID         uint64
	SubDstAppID      uint64
	AppClientVersion uint32
	Uin              uint64
	I2               uint16
	IPv4Address      net.IP
	CurrentTime      uint32
	PasswordMD5      [16]byte
	TGTGTKey         [16]byte
	LoginType        uint32
	T106             []byte
	T16A             []byte
	MiscBitmap       uint32
	SubSigMap        uint32
	SubAppIDList     []uint64
	MainSigMap       uint32
	SrcAppID         uint64
	I7               uint16
	I8               uint8
	I9               uint16
	I10              uint8
	KSID             []byte
	T104             []byte
	PackageName      []byte
	Domains          []string
}

func NewAuthGetSessionTicketWithPasswordRequest(uin uint64, password string) *AuthGetSessionTicketWithPasswordRequest {
	return &AuthGetSessionTicketWithPasswordRequest{
		Username:  fmt.Sprintf("%d", uin),
		ImageType: 0x01,

		DstAppID:         defaultClientDstAppID,
		SubDstAppID:      defaultClientOpenAppID,
		AppClientVersion: 0x00000000,
		Uin:              uin,
		I2:               0x0000,
		IPv4Address:      defaultDeviceIPv4Address,
		CurrentTime:      uint32(time.Now().UnixNano() / 1e6), // TODO: sync server time
		PasswordMD5:      md5.Sum([]byte(password)),
		TGTGTKey:         [16]byte{},
		LoginType:        0x00000001,
		T106:             nil,
		T16A:             nil,
		MiscBitmap:       defaultClientMiscBitmap,
		SubSigMap:        defaultClientSubSigMap,
		SubAppIDList:     defaultClientSubAppIDList,
		MainSigMap:       defaultClientMainSigMap & 0xfdfffffe,
		SrcAppID:         defaultClientOpenAppID,
		I7:               0x0000,
		I8:               0x00,
		I9:               0x0000,
		I10:              0x01,
		KSID:             defaultDeviceKSID,
		T104:             nil,
		PackageName:      defaultClientPackageName,
		Domains:          defaultClientDomains,
	}
}

func (req *AuthGetSessionTicketWithPasswordRequest) EncodeOICQMessage(ctx context.Context) (*message.OICQMessage, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0018] = tlv.NewT18(req.DstAppID, req.AppClientVersion, req.Uin, req.I2)
	tlvs[0x0001] = tlv.NewT1(req.Uin, req.IPv4Address)
	if len(req.T106) == 0 {
		tlvs[0x0106] = tlv.NewT106(req.DstAppID, req.SubDstAppID, req.AppClientVersion, req.Uin, req.CurrentTime, req.IPv4Address, true, req.PasswordMD5, 0, req.Username, req.TGTGTKey, true, defaultDeviceGUID, req.LoginType)
	} else {
		tlvs[0x0106] = tlv.NewTLV(0x0106, 0x0000, bytes.NewBuffer(req.T106))
	}
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0100] = tlv.NewT100(req.DstAppID, req.SrcAppID, req.AppClientVersion, req.MainSigMap)
	tlvs[0x0107] = tlv.NewT107(req.I7, req.I8, req.I9, req.I10)
	if len(req.KSID) != 0 {
		tlvs[0x0108] = tlv.NewT108(req.KSID)
	}
	if len(req.T104) != 0 {
		tlvs[0x0104] = tlv.NewT104(req.T104)
	}
	tlvs[0x0142] = tlv.NewT142(req.PackageName)
	if !util.CheckUinAccount(req.Username) {
		tlvs[0x0112] = tlv.NewT112([]byte(req.Username))
	}
	tlvs[0x0144] = tlv.NewT144(req.TGTGTKey,
		tlv.NewT109(md5.Sum(defaultDeviceOSID)),
		tlv.NewT52D(ctx),
		tlv.NewT124(defaultDeviceOSType, defaultDeviceOSVersion, defaultDeviceNetworkTypeID, defaultDeviceSIMOPName, nil, defaultDeviceAPNName),
		tlv.NewT128(defaultDeviceIsGUIDFileNil, defaultDeviceIsGUIDGenSucc, defaultDeviceIsGUIDChanged, defaultDeviceGUIDFlag, defaultDeviceOSBuildModel, defaultDeviceGUID, defaultDeviceOSBuildBrand),
		tlv.NewT16E(defaultDeviceOSBuildModel),
	)
	tlvs[0x0145] = tlv.NewT145(defaultDeviceGUID)
	tlvs[0x0147] = tlv.NewT147(req.DstAppID, defaultClientVersionName, defaultClientSignatureMD5)
	if req.MiscBitmap&0x80 != 0 {
		tlvs[0x0166] = tlv.NewT166(req.ImageType)
	}
	if len(req.T16A) != 0 {
		tlvs[0x016a] = tlv.NewT16A(req.T16A)
	}
	tlvs[0x0154] = tlv.NewT154(req.Seq)
	tlvs[0x0141] = tlv.NewT141(defaultDeviceSIMOPName, defaultDeviceNetworkTypeID, defaultDeviceAPNName)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	if len(req.Domains) > 0 {
		tlvs[0x0511] = tlv.NewT511(req.Domains)
	}
	if len(req.T172) != 0 {
		tlvs[0x0172] = tlv.NewT172(req.T172)
	}
	if req.LoginType == 0x000000003 {
		tlvs[0x0185] = tlv.NewT185(0x01)
	}
	if false {
		tlvs[0x0400] = tlv.NewT400([16]byte{}, req.Uin, nil, [16]byte{}, req.DstAppID, req.SubDstAppID, nil)
	}
	tlvs[0x0187] = tlv.NewT187(md5.Sum(defaultDeviceMACAddress))
	tlvs[0x0188] = tlv.NewT188(md5.Sum(defaultDeviceOSID))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(defaultDeviceIMSI)))
	tlvs[0x0191] = tlv.NewT191(defaultClientVerifyMethod)
	// DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(md5.Sum(defaultDeviceBSSIDAddress), defaultDeviceSSIDAddress)
	tlvs[0x0177] = tlv.NewT177(defaultClientBuildTime, defaultClientSDKVersion)
	tlvs[0x0516] = tlv.NewTLV(0x0516, 0x0004, bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00}))
	tlvs[0x0521] = tlv.NewTLV(0x0521, 0x0006, bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
	if len(req.T9) != 0 {
		buf := bytes.NewBuffer([]byte{})
		tlv.NewTLV(0x0536, 0x0002, bytes.NewBuffer(req.T9)).Encode(buf)
		tlvs[0x0525] = tlv.NewTLV(0x0525, 0x0000, buf)
	}
	if len(req.TGTQR) != 0 {
		tlvs[0x0318] = tlv.NewTLV(0x0318, 0x0000, bytes.NewBuffer(req.TGTQR))
	}
	// DISABLED: tgt
	// tlvs[0x0544] = tlv.NewT544(req.Username, "810_9", nil)
	// DISABLED: tgtgt qimei
	// tlvs[0x0545] = tlv.NewT545(md5.Sum([]byte("qimei")))
	// DISABLED: nativeGetTestData
	// tlvs[0x0548] = tlv.NewT548([]byte("nativeGetTestData"))

	return &message.OICQMessage{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: 0x07,
		RandomKey:     defaultClientRandomKey,
		PublicKey:     ecdh.PublicKey,
		ShareKey:      ecdh.ShareKey,
		Type:          0x0009,
		TLVs:          tlvs,
	}, nil
}

func (req *AuthGetSessionTicketWithPasswordRequest) Encode(ctx context.Context) (*ClientToServerMessage, error) {
	msg, err := req.EncodeOICQMessage(ctx)
	if err != nil {
		return nil, err
	}
	buf, err := message.MarshalOICQMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &ClientToServerMessage{
		Seq:      req.Seq,
		Username: req.Username,
		Buffer:   buf,
		Simple:   false,
	}, nil
}

func (c *Client) AuthGetSessionTicketWithPassword(ctx context.Context, req *AuthGetSessionTicketWithPasswordRequest) (interface{}, error) {
	req.Seq = c.getNextSeq()
	req.TGTGTKey = [16]byte{}
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
