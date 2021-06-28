package rpc

import (
	"context"
	"crypto/md5"
	"net"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/tlv"
	"github.com/elap5e/go-mobileqq-api/util"
)

type AuthGetSessionTicketsWithPasswordRequest struct {
	Seq            uint32
	HashGUID       [16]byte
	RandomPassword [16]byte
	RandomSeed     []byte
	LoginExtraData []byte
	T172           []byte
	TGTQR          []byte

	Username string

	DstAppID         uint64
	SubDstAppID      uint64
	AppClientVersion uint32 // constant 0x00000000
	Uin              uint64
	I2               uint16 // constant 0x0000
	IPv4Address      net.IP
	ServerTime       uint32
	PasswordMD5      [16]byte
	TGTGTKey         [16]byte // placeholder
	LoginType        uint32   // 0x00, 0x01, 0x03
	T106             []byte
	T16A             []byte
	MiscBitmap       uint32
	SubSigMap        uint32
	SubAppIDList     []uint64
	MainSigMap       uint32
	SrcAppID         uint64
	I7               uint16 // constant 0x0000
	I8               uint8  // constant 0x00
	I9               uint16 // constant 0x0000
	I10              uint8  // constant 0x01
	KSID             []byte
	T104             []byte // placeholder
	PackageName      []byte
	Domains          []string
}

func NewAuthGetSessionTicketsWithPasswordRequest(username string, password string) *AuthGetSessionTicketsWithPasswordRequest {
	uin, _ := strconv.ParseInt(username, 10, 64)
	return &AuthGetSessionTicketsWithPasswordRequest{
		Username: username,

		DstAppID:         defaultClientDstAppID,
		SubDstAppID:      defaultClientOpenAppID,
		AppClientVersion: 0x00000000,
		Uin:              uint64(uin),
		I2:               0x0000,
		IPv4Address:      defaultDeviceIPv4Address,
		ServerTime:       util.GetServerTime(),
		PasswordMD5:      md5.Sum([]byte(password)),
		TGTGTKey:         [16]byte{},
		LoginType:        0x00000001,
		T106:             nil, // nil
		T16A:             nil, // nil
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
	}
}

func (req *AuthGetSessionTicketsWithPasswordRequest) GetTLVs(ctx context.Context) (map[uint16]tlv.TLVCodec, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0018] = tlv.NewT18(req.DstAppID, req.AppClientVersion, req.Uin, req.I2)
	tlvs[0x0001] = tlv.NewT1(req.Uin, req.IPv4Address)
	if len(req.T106) == 0 {
		tlvs[0x0106] = tlv.NewT106(req.DstAppID, req.SubDstAppID, req.AppClientVersion, req.Uin, req.ServerTime, req.IPv4Address, true, req.PasswordMD5, 0, req.Username, req.TGTGTKey, true, deviceGUID[:], req.LoginType)
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
	if !util.CheckUsername(req.Username) {
		tlvs[0x0112] = tlv.NewT112([]byte(req.Username))
	}
	tlvs[0x0144] = tlv.NewT144(req.TGTGTKey,
		tlv.NewT109(md5.Sum(defaultDeviceOSBuildID)),
		tlv.NewT52D(ctx),
		tlv.NewT124([]byte(defaultDeviceOSType), []byte(defaultDeviceOSVersion), defaultDeviceNetworkTypeID, defaultDeviceSIMOPName, nil, defaultDeviceAPNName),
		tlv.NewT128(deviceIsGUIDFileNil, deviceIsGUIDGenSucc, deviceIsGUIDChanged, deviceGUIDFlag, []byte(defaultDeviceOSBuildModel), deviceGUID[:], defaultDeviceOSBuildBrand),
		tlv.NewT16E([]byte(defaultDeviceOSBuildModel)),
	)
	tlvs[0x0145] = tlv.NewT145(deviceGUID[:])
	tlvs[0x0147] = tlv.NewT147(req.DstAppID, clientVersionName, clientSignatureMD5)
	if req.MiscBitmap&0x80 != 0 {
		tlvs[0x0166] = tlv.NewT166(clientImageType)
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
	if false { // TODO: code2d
		tlvs[0x0400] = tlv.NewT400(req.HashGUID, req.Uin, deviceGUID, req.RandomPassword, req.DstAppID, req.SubDstAppID, req.RandomSeed)
	}
	tlvs[0x0187] = tlv.NewT187(md5.Sum(defaultDeviceMACAddress))
	tlvs[0x0188] = tlv.NewT188(md5.Sum(defaultDeviceOSBuildID))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(defaultDeviceIMSI)))
	tlvs[0x0191] = tlv.NewT191(clientVerifyMethod)
	// DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(md5.Sum(defaultDeviceBSSIDAddress), defaultDeviceSSIDAddress)
	tlvs[0x0177] = tlv.NewT177(clientBuildTime, clientSDKVersion)
	tlvs[0x0516] = tlv.NewTLV(0x0516, 0x0004, bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00}))             // SourceType
	tlvs[0x0521] = tlv.NewTLV(0x0521, 0x0006, bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00})) // ProductType
	if len(req.LoginExtraData) != 0 {
		buf := bytes.NewBuffer([]byte{})
		buf.EncodeUint16(0x0001)
		tlv.NewTLV(0x0536, 0x0002, bytes.NewBuffer(req.LoginExtraData)).Encode(buf)
		tlvs[0x0525] = tlv.NewTLV(0x0525, 0x0000, buf)
	}
	if len(req.TGTQR) != 0 { // TODO: code2d
		tlvs[0x0318] = tlv.NewTLV(0x0318, 0x0000, bytes.NewBuffer(req.TGTQR))
	}
	// DISABLED: tgt
	// tlvs[0x0544] = tlv.NewT544(req.Uin, deviceGUID, clientSDKVersion, 0x0009)
	// DISABLED: tgtgt qimei
	// tlvs[0x0545] = tlv.NewT545(md5.Sum([]byte("qimei")))
	// DISABLED: nativeGetTestData
	// tlvs[0x0548] = tlv.NewT548([]byte("nativeGetTestData"))
	return tlvs, nil
}

func (c *Client) AuthGetSessionTicketsWithPassword(ctx context.Context, req *AuthGetSessionTicketsWithPasswordRequest) (*AuthGetSessionTicketsResponse, error) {
	req.Seq = c.getNextSeq()
	req.HashGUID = c.hashGUID
	req.RandomPassword = c.randomPassword
	req.LoginExtraData = c.loginExtraData
	req.T172 = c.t172
	req.TGTGTKey = c.tgtgtKey
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
