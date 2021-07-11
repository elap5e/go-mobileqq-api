package auth

import (
	"context"
	"crypto/md5"
	"net"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/tlv"
	"github.com/elap5e/go-mobileqq-api/util"
)

type getSessionTicketsWithPasswordRequest struct {
	request

	DstAppID         uint64
	SubDstAppID      uint64
	AppClientVersion uint32 // constant 0x00000000
	_Uin             uint64
	I2               uint16 // constant 0x0000
	_IPv4Address     net.IP // h.opt.Client.MiscBitmap
	ServerTime       uint32
	PasswordMD5      [16]byte
	_UserA1Key       [16]byte // c.userA1Key
	LoginType        uint32   // 0x00, 0x01, 0x03
	UserA1           []byte
	T16A             []byte
	_MiscBitmap      uint32 // h.opt.Client.MiscBitmap
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
	_PackageName     []byte // []byte(h.opt.Client.PackageName)
	Domains          []string
}

func newGetSessionTicketsWithPasswordRequest(
	username string,
	password string,
) *getSessionTicketsWithPasswordRequest {
	req := &getSessionTicketsWithPasswordRequest{
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

func (req *getSessionTicketsWithPasswordRequest) MustGetTLVs(
	ctx context.Context,
) map[uint16]tlv.TLVCodec {
	h := ForHandler(ctx)
	sig := h.client.GetUserSignature(req.GetUsername())
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0018] = tlv.NewT18(
		req.DstAppID,
		req.AppClientVersion,
		req.GetUin(),
		req.I2,
	)
	tlvs[0x0001] = tlv.NewT1(req.GetUin(), h.opt.Device.IPv4Address)
	if len(sig.Tickets["A1"].Sig) == 0 {
		tlvs[0x0106] = tlv.NewT106(
			req.DstAppID,
			req.SubDstAppID,
			req.AppClientVersion,
			req.GetUin(),
			req.ServerTime,
			h.opt.Device.IPv4Address,
			true,
			req.PasswordMD5,
			0,
			req.GetUsername(),
			util.BytesToSTBytes(sig.Tickets["A1"].Key),
			true,
			h.opt.Device.GUID,
			req.LoginType,
			h.opt.Client.SSOVersion,
		)
	} else {
		tlvs[0x0106] = tlv.NewTLV(
			0x0106,
			0x0000,
			bytes.NewBuffer(sig.Tickets["A1"].Sig),
		)
	}
	tlvs[0x0116] = tlv.NewT116(
		h.opt.Client.MiscBitmap,
		req.SubSigMap,
		req.SubAppIDList,
	)
	tlvs[0x0100] = tlv.NewT100(
		req.DstAppID,
		req.SrcAppID,
		req.AppClientVersion,
		req.MainSigMap,
		h.opt.Client.SSOVersion,
	)
	tlvs[0x0107] = tlv.NewT107(req.I7, req.I8, req.I9, req.I10)
	if len(sig.Session.KSID) != 0 {
		tlvs[0x0108] = tlv.NewT108(sig.Session.KSID)
	}
	if len(sig.Session.Auth) != 0 {
		tlvs[0x0104] = tlv.NewT104(sig.Session.Auth)
	}
	tlvs[0x0142] = tlv.NewT142([]byte(h.opt.Client.PackageName))
	if !util.CheckUsername(req.GetUsername()) {
		tlvs[0x0112] = tlv.NewT112([]byte(req.GetUsername()))
	}
	tlvs[0x0144] = tlv.NewT144(util.BytesToSTBytes(sig.Tickets["A1"].Key),
		tlv.NewT109(md5.Sum([]byte(h.opt.Device.OSBuildID))),
		tlv.NewT52D(&pb.DeviceReport{
			Bootloader:   []byte(h.opt.Device.Bootloader),
			ProcVersion:  []byte(h.opt.Device.ProcVersion),
			Codename:     []byte(h.opt.Device.Codename),
			Incremental:  []byte(h.opt.Device.Incremental),
			Fingerprint:  []byte(h.opt.Device.Fingerprint),
			BootId:       []byte(h.opt.Device.BootID),
			AndroidId:    []byte(h.opt.Device.OSBuildID),
			Baseband:     []byte(h.opt.Device.Baseband),
			InnerVersion: []byte(h.opt.Device.InnerVersion),
		}),
		tlv.NewT124(
			[]byte(defaultDeviceOSType),
			[]byte(defaultDeviceOSVersion),
			h.opt.Device.NetworkType,
			defaultDeviceSIMOPName,
			nil,
			defaultDeviceAPNName,
		),
		tlv.NewT128(
			h.opt.Device.IsGUIDFileNil,
			h.opt.Device.IsGUIDGenSucc,
			h.opt.Device.IsGUIDChanged,
			h.opt.Device.GUIDFlag,
			[]byte(defaultDeviceOSBuildModel),
			h.opt.Device.GUID,
			defaultDeviceOSBuildBrand,
		),
		tlv.NewT16E([]byte(defaultDeviceOSBuildModel)),
	)
	tlvs[0x0145] = tlv.NewT145(util.BytesToSTBytes(h.opt.Device.GUID))
	tlvs[0x0147] = tlv.NewT147(
		req.DstAppID,
		[]byte(h.opt.Client.VersionName),
		util.BytesToSTBytes(h.opt.Client.SignatureMD5),
	)
	if h.opt.Client.MiscBitmap&0x80 != 0 {
		tlvs[0x0166] = tlv.NewT166(h.opt.Client.ImageType)
	}
	if len(h.t16a) != 0 {
		tlvs[0x016a] = tlv.NewT16A(h.t16a)
	}
	tlvs[0x0154] = tlv.NewT154(req.GetSeq())
	tlvs[0x0141] = tlv.NewT141(
		defaultDeviceSIMOPName,
		h.opt.Device.NetworkType,
		defaultDeviceAPNName,
	)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	if len(req.Domains) > 0 {
		tlvs[0x0511] = tlv.NewT511(req.Domains)
	}
	if len(h.t172) != 0 {
		tlvs[0x0172] = tlv.NewT172(h.t172)
	}
	if req.LoginType == 0x000000003 {
		tlvs[0x0185] = tlv.NewT185(0x01)
	}
	// if false { // TODO: code2d
	// 	tlvs[0x0400] = tlv.NewT400(
	// 		h.hashedGUID,
	// 		req.GetUin(),
	// 		util.BytesToSTBytes(h.opt.Device.GUID),
	// 		h.randomPassword,
	// 		req.DstAppID,
	// 		req.SubDstAppID,
	// 		h.randomSeed,
	// 	)
	// }
	tlvs[0x0187] = tlv.NewT187(md5.Sum([]byte(h.opt.Device.MACAddress)))
	tlvs[0x0188] = tlv.NewT188(md5.Sum([]byte(h.opt.Device.OSBuildID)))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(h.opt.Device.IMSI)))
	if h.opt.Client.CanCaptcha {
		tlvs[0x0191] = tlv.NewT191(0x82)
	} else {
		tlvs[0x0191] = tlv.NewT191(0x00)
	}
	// // DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(
		md5.Sum([]byte(h.opt.Device.BSSIDAddress)),
		[]byte(h.opt.Device.SSIDAddress),
	)
	tlvs[0x0177] = tlv.NewT177(h.opt.Client.BuildTime, h.opt.Client.SDKVersion)
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
	if len(h.loginExtraData) != 0 {
		buf := bytes.NewBuffer([]byte{})
		buf.EncodeUint16(0x0001)
		tlv.NewTLV(
			0x0536,
			0x0002,
			bytes.NewBuffer(h.loginExtraData),
		).Encode(buf)
		tlvs[0x0525] = tlv.NewTLV(0x0525, 0x0000, buf)
	}
	// if len(h.tgtQR) != 0 { // TODO: code2d
	// 	tlvs[0x0318] = tlv.NewTLV(0x0318, 0x0000, bytes.NewBuffer(h.tgtQR))
	// }
	// // DISABLED: tgt
	// tlvs[0x0544] = tlv.NewT544(
	// 	req.Uin,
	// 	h.opt.Device.GUID,
	// 	h.opt.Client.SDKVersion,
	// 	0x0009,
	// )
	// // DISABLED: tgtgt qimei
	// tlvs[0x0545] = tlv.NewT545(md5.Sum([]byte("qimei")))
	// // DISABLED: nativeGetTestData
	// tlvs[0x0548] = tlv.NewT548([]byte("nativeGetTestData"))
	req.SetType(0x0009)
	req.SetServiceMethod(ServiceMethodAuthLogin)
	return tlvs
}

func (h *Handler) getSessionTicketsWithPassword(
	ctx context.Context,
	req *getSessionTicketsWithPasswordRequest,
) (*Response, error) {
	return h.getSessionTickets(ctx, req)
}
