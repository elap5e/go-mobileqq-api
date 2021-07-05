package client

import (
	"context"
	"crypto/md5"
	"net"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/tlv"
	"github.com/elap5e/go-mobileqq-api/util"
)

type AuthGetSessionTicketsWithPasswordRequest struct {
	authGetSessionTicketsRequest

	DstAppID         uint64
	SubDstAppID      uint64
	AppClientVersion uint32 // constant 0x00000000
	_Uin             uint64
	I2               uint16 // constant 0x0000
	_IPv4Address     net.IP // c.cfg.Client.MiscBitmap
	ServerTime       uint32
	PasswordMD5      [16]byte
	_UserA1Key       [16]byte // c.userA1Key
	LoginType        uint32   // 0x00, 0x01, 0x03
	UserA1           []byte
	T16A             []byte
	_MiscBitmap      uint32 // c.cfg.Client.MiscBitmap
	SubSigMap        uint32
	SubAppIDList     []uint64
	MainSigMap       uint32
	SrcAppID         uint64
	I7               uint16 // constant 0x0000
	I8               uint8  // constant 0x00
	I9               uint16 // constant 0x0000
	I10              uint8  // constant 0x01
	_KSID            []byte // sig.Session.KSID
	_AuthSession     []byte // sig.Session.AuthSession
	_PackageName     []byte // []byte(c.cfg.Client.PackageName)
	Domains          []string
}

func NewAuthGetSessionTicketsWithPasswordRequest(
	username string,
	password string,
) *AuthGetSessionTicketsWithPasswordRequest {
	req := &AuthGetSessionTicketsWithPasswordRequest{
		DstAppID:         defaultClientDstAppID,
		SubDstAppID:      defaultClientOpenAppID,
		AppClientVersion: 0x00000000,
		_Uin:             0x00000000,
		I2:               0x0000,
		_IPv4Address:     nil,
		ServerTime:       util.GetServerTime(),
		PasswordMD5:      md5.Sum([]byte(password)),
		_UserA1Key:       [16]byte{},
		LoginType:        0x00000001,
		UserA1:           nil, // nil
		T16A:             nil, // nil
		_MiscBitmap:      0x00000000,
		SubSigMap:        defaultClientSubSigMap,
		SubAppIDList:     defaultClientSubAppIDList,
		MainSigMap:       defaultClientMainSigMap & 0xfdfffffe,
		SrcAppID:         defaultClientOpenAppID,
		I7:               0x0000,
		I8:               0x00,
		I9:               0x0000,
		I10:              0x01,
		_KSID:            nil,
		_AuthSession:     nil,
		_PackageName:     []byte{},
		Domains:          defaultClientDomains,
	}
	req.SetUsername(username)
	return req
}

func (req *AuthGetSessionTicketsWithPasswordRequest) GetTLVs(
	ctx context.Context,
) (map[uint16]tlv.TLVCodec, error) {
	c := ForClient(ctx)
	sig := c.rpc.GetUserSignature(req.GetUsername())
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0018] = tlv.NewT18(
		req.DstAppID,
		req.AppClientVersion,
		req.GetUin(),
		req.I2,
	)
	tlvs[0x0001] = tlv.NewT1(req.GetUin(), c.cfg.Device.IPv4Address)
	if len(sig.Tickets["A1"].Sig) == 0 {
		tlvs[0x0106] = tlv.NewT106(
			req.DstAppID,
			req.SubDstAppID,
			req.AppClientVersion,
			req.GetUin(),
			req.ServerTime,
			c.cfg.Device.IPv4Address,
			true,
			req.PasswordMD5,
			0,
			req.GetUsername(),
			util.BytesToSTBytes(sig.Tickets["A1"].Key),
			true,
			c.cfg.Device.GUID,
			req.LoginType,
			c.cfg.Client.SSOVersion,
		)
	} else {
		tlvs[0x0106] = tlv.NewTLV(
			0x0106,
			0x0000,
			bytes.NewBuffer(sig.Tickets["A1"].Sig),
		)
	}
	tlvs[0x0116] = tlv.NewT116(
		c.cfg.Client.MiscBitmap,
		req.SubSigMap,
		req.SubAppIDList,
	)
	tlvs[0x0100] = tlv.NewT100(
		req.DstAppID,
		req.SrcAppID,
		req.AppClientVersion,
		req.MainSigMap,
		c.cfg.Client.SSOVersion,
	)
	tlvs[0x0107] = tlv.NewT107(req.I7, req.I8, req.I9, req.I10)
	if len(sig.Session.KSID) != 0 {
		tlvs[0x0108] = tlv.NewT108(sig.Session.KSID)
	}
	if len(sig.Session.Auth) != 0 {
		tlvs[0x0104] = tlv.NewT104(sig.Session.Auth)
	}
	tlvs[0x0142] = tlv.NewT142([]byte(c.cfg.Client.PackageName))
	if !util.CheckUsername(req.GetUsername()) {
		tlvs[0x0112] = tlv.NewT112([]byte(req.GetUsername()))
	}
	tlvs[0x0144] = tlv.NewT144(util.BytesToSTBytes(sig.Tickets["A1"].Key),
		tlv.NewT109(md5.Sum([]byte(c.cfg.Device.OSBuildID))),
		tlv.NewT52D(&pb.DeviceReport{
			Bootloader:   []byte(c.cfg.Device.Bootloader),
			ProcVersion:  []byte(c.cfg.Device.ProcVersion),
			Codename:     []byte(c.cfg.Device.Codename),
			Incremental:  []byte(c.cfg.Device.Incremental),
			Fingerprint:  []byte(c.cfg.Device.Fingerprint),
			BootId:       []byte(c.cfg.Device.BootID),
			AndroidId:    []byte(c.cfg.Device.OSBuildID),
			Baseband:     []byte(c.cfg.Device.Baseband),
			InnerVersion: []byte(c.cfg.Device.InnerVersion),
		}),
		tlv.NewT124(
			[]byte(defaultDeviceOSType),
			[]byte(defaultDeviceOSVersion),
			c.cfg.Device.NetworkType,
			defaultDeviceSIMOPName,
			nil,
			defaultDeviceAPNName,
		),
		tlv.NewT128(
			c.cfg.Device.IsGUIDFileNil,
			c.cfg.Device.IsGUIDGenSucc,
			c.cfg.Device.IsGUIDChanged,
			c.cfg.Device.GUIDFlag,
			[]byte(defaultDeviceOSBuildModel),
			c.cfg.Device.GUID,
			defaultDeviceOSBuildBrand,
		),
		tlv.NewT16E([]byte(defaultDeviceOSBuildModel)),
	)
	tlvs[0x0145] = tlv.NewT145(util.BytesToSTBytes(c.cfg.Device.GUID))
	tlvs[0x0147] = tlv.NewT147(
		req.DstAppID,
		[]byte(c.cfg.Client.VersionName),
		util.BytesToSTBytes(c.cfg.Client.SignatureMD5),
	)
	if c.cfg.Client.MiscBitmap&0x80 != 0 {
		tlvs[0x0166] = tlv.NewT166(c.cfg.Client.ImageType)
	}
	if len(c.t16a) != 0 {
		tlvs[0x016a] = tlv.NewT16A(c.t16a)
	}
	tlvs[0x0154] = tlv.NewT154(req.GetSeq())
	tlvs[0x0141] = tlv.NewT141(
		defaultDeviceSIMOPName,
		c.cfg.Device.NetworkType,
		defaultDeviceAPNName,
	)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	if len(req.Domains) > 0 {
		tlvs[0x0511] = tlv.NewT511(req.Domains)
	}
	if len(c.t172) != 0 {
		tlvs[0x0172] = tlv.NewT172(c.t172)
	}
	if req.LoginType == 0x000000003 {
		tlvs[0x0185] = tlv.NewT185(0x01)
	}
	// if false { // TODO: code2d
	// 	tlvs[0x0400] = tlv.NewT400(
	// 		c.hashedGUID,
	// 		req.GetUin(),
	// 		util.BytesToSTBytes(c.cfg.Device.GUID),
	// 		c.randomPassword,
	// 		req.DstAppID,
	// 		req.SubDstAppID,
	// 		c.randomSeed,
	// 	)
	// }
	tlvs[0x0187] = tlv.NewT187(md5.Sum([]byte(c.cfg.Device.MACAddress)))
	tlvs[0x0188] = tlv.NewT188(md5.Sum([]byte(c.cfg.Device.OSBuildID)))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(c.cfg.Device.IMSI)))
	if c.cfg.Client.CanCaptcha {
		tlvs[0x0191] = tlv.NewT191(0x82)
	} else {
		tlvs[0x0191] = tlv.NewT191(0x00)
	}
	// // DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(
		md5.Sum([]byte(c.cfg.Device.BSSIDAddress)),
		[]byte(c.cfg.Device.SSIDAddress),
	)
	tlvs[0x0177] = tlv.NewT177(c.cfg.Client.BuildTime, c.cfg.Client.SDKVersion)
	tlvs[0x0516] = tlv.NewTLV(
		0x0516,
		0x0004,
		bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00}),
	) // SourceType
	tlvs[0x0521] = tlv.NewTLV(
		0x0521,
		0x0006,
		bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
	) // ProductType
	if len(c.loginExtraData) != 0 {
		buf := bytes.NewBuffer([]byte{})
		buf.EncodeUint16(0x0001)
		tlv.NewTLV(
			0x0536,
			0x0002,
			bytes.NewBuffer(c.loginExtraData),
		).Encode(buf)
		tlvs[0x0525] = tlv.NewTLV(0x0525, 0x0000, buf)
	}
	// if len(c.tgtQR) != 0 { // TODO: code2d
	// 	tlvs[0x0318] = tlv.NewTLV(0x0318, 0x0000, bytes.NewBuffer(c.tgtQR))
	// }
	// // DISABLED: tgt
	// tlvs[0x0544] = tlv.NewT544(
	// 	req.Uin,
	// 	c.cfg.Device.GUID,
	// 	c.cfg.Client.SDKVersion,
	// 	0x0009,
	// )
	// // DISABLED: tgtgt qimei
	// tlvs[0x0545] = tlv.NewT545(md5.Sum([]byte("qimei")))
	// // DISABLED: nativeGetTestData
	// tlvs[0x0548] = tlv.NewT548([]byte("nativeGetTestData"))
	req.SetType(0x0009)
	req.SetServiceMethod(ServiceMethodAuthLogin)
	return tlvs, nil
}

func (c *Client) AuthGetSessionTicketsWithPassword(
	ctx context.Context,
	req *AuthGetSessionTicketsWithPasswordRequest,
) (*AuthGetSessionTicketsResponse, error) {
	return c.AuthGetSessionTickets(ctx, req)
}
