package rpc

import (
	"context"
	"crypto/md5"

	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/tlv"
	"github.com/elap5e/go-mobileqq-api/util"
)

type AuthGetSessionTicketsWithoutPasswordRequest struct {
	authGetSessionTicketsRequest

	_Uin             uint64
	DstAppID         uint64
	SrcAppID         uint64 // constant 0x00000064
	AppClientVersion uint32
	MainSigMap       uint32
	_UserA2          []byte // sig.Tickets["A2"].Sig
	_MiscBitmap      uint32 // c.cfg.Client.MiscBitmap
	SubSigMap        uint32
	SubAppIDList     []uint64
	_KSID            []byte // sig.Session.KSID
	_UserD2          []byte // sig.Tickets["D2"].Sig
	Domains          []string

	changeD2 bool
}

func NewAuthGetSessionTicketsWithoutPasswordRequest(
	username string, changeD2 ...bool,
) *AuthGetSessionTicketsWithoutPasswordRequest {
	req := &AuthGetSessionTicketsWithoutPasswordRequest{
		_Uin:             0x00000000,
		DstAppID:         defaultClientDstAppID,
		SrcAppID:         0x00000064,
		AppClientVersion: 0x00000000,
		MainSigMap:       defaultClientMainSigMap & 0xfdfffffe,
		_UserA2:          nil,
		_MiscBitmap:      0x00000000,
		SubSigMap:        defaultClientSubSigMap,
		SubAppIDList:     defaultClientSubAppIDList,
		_KSID:            nil,
		_UserD2:          nil,
		Domains:          defaultClientDomains,
	}
	if len(changeD2) > 0 {
		req.changeD2 = changeD2[0]
	} else {
		req.changeD2 = true
	}
	req.SetUsername(username)
	return req
}

func (req *AuthGetSessionTicketsWithoutPasswordRequest) GetTLVs(
	ctx context.Context,
) (map[uint16]tlv.TLVCodec, error) {
	c := ForClient(ctx)
	sig := c.GetUserSignature(req.GetUsername())
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0100] = tlv.NewT100(
		req.DstAppID,
		req.SrcAppID,
		req.AppClientVersion,
		req.MainSigMap,
		c.cfg.Client.SSOVersion,
	)
	tlvs[0x010a] = tlv.NewT10A(sig.Tickets["A2"].Sig)
	tlvs[0x0116] = tlv.NewT116(c.cfg.Client.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0108] = tlv.NewT108(sig.Session.KSID)
	key := [16]byte{}
	if !req.changeD2 {
		copy(key[:], sig.Tickets["A2"].Key)
	} else {
		key = md5.Sum(sig.Tickets["D2"].Key)
	}
	tlvs[0x0144] = tlv.NewT144(key,
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
	if !req.changeD2 {
		tlvs[0x0145] = tlv.NewT145(util.BytesToSTBytes(c.cfg.Device.GUID))
	} else {
		tlvs[0x0143] = tlv.NewT143(sig.Tickets["D2"].Sig)
	}
	tlvs[0x0142] = tlv.NewT142([]byte(c.cfg.Client.PackageName))
	tlvs[0x0154] = tlv.NewT154(req.GetSeq())
	tlvs[0x0018] = tlv.NewT18(
		req.DstAppID,
		req.MainSigMap,
		req.GetUin(),
		0x0000,
	)
	tlvs[0x0141] = tlv.NewT141(
		defaultDeviceSIMOPName,
		c.cfg.Device.NetworkType,
		defaultDeviceAPNName,
	)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	if len(req.Domains) > 0 {
		tlvs[0x0511] = tlv.NewT511(req.Domains)
	}
	tlvs[0x0147] = tlv.NewT147(
		req.DstAppID,
		[]byte(c.cfg.Client.VersionName),
		util.BytesToSTBytes(c.cfg.Client.SignatureMD5),
	)
	if len(c.t172) != 0 {
		tlvs[0x0172] = tlv.NewT172(c.t172)
	}
	tlvs[0x0177] = tlv.NewT177(c.cfg.Client.BuildTime, c.cfg.Client.SDKVersion)
	tlvs[0x0187] = tlv.NewT187(md5.Sum([]byte(c.cfg.Device.MACAddress)))
	tlvs[0x0188] = tlv.NewT188(md5.Sum([]byte(c.cfg.Device.OSBuildID)))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(c.cfg.Device.IMSI)))
	// // DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(
		md5.Sum([]byte(c.cfg.Device.BSSIDAddress)),
		[]byte(c.cfg.Device.SSIDAddress),
	)
	// // DISABLED: tgt
	// tlvs[0x0544] = tlv.NewT544(
	// 	req.Uin,
	// 	c.cfg.Device.GUID,
	// 	c.cfg.Client.SDKVersion,
	// 	0x0009,
	// )
	if !req.changeD2 {
		req.SetType(0x000a)
	} else {
		req.SetType(0x000b)
	}
	req.SetServiceMethod(ServiceMethodAuthExchangeAccount)
	return tlvs, nil
}

func (c *Client) AuthGetSessionTicketsWithoutPassword(
	ctx context.Context,
	req *AuthGetSessionTicketsWithoutPasswordRequest,
) (*AuthGetSessionTicketsResponse, error) {
	return c.AuthGetSessionTickets(ctx, req)
}
