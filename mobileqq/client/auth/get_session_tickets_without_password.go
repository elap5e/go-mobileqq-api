package auth

import (
	"context"
	"crypto/md5"

	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/tlv"
	"github.com/elap5e/go-mobileqq-api/util"
)

type getSessionTicketsWithoutPasswordRequest struct {
	request

	_Uin             uint64
	DstAppID         uint64
	SrcAppID         uint64 // constant 0x00000064
	AppClientVersion uint32
	MainSigMap       uint32
	_UserA2          []byte // sig.Tickets["A2"].Sig
	_MiscBitmap      uint32 // h.opt.Client.MiscBitmap
	SubSigMap        uint32
	SubAppIDList     []uint64
	_KSID            []byte // sig.Session.KSID
	_UserD2          []byte // sig.Tickets["D2"].Sig
	Domains          []string

	changeD2 bool
}

func newGetSessionTicketsWithoutPasswordRequest(
	username string, changeD2 ...bool,
) *getSessionTicketsWithoutPasswordRequest {
	req := &getSessionTicketsWithoutPasswordRequest{
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

func (req *getSessionTicketsWithoutPasswordRequest) MustGetTLVs(
	ctx context.Context,
) map[uint16]tlv.TLVCodec {
	h := ForHandler(ctx)
	sig := h.client.GetUserSignature(req.GetUsername())
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0100] = tlv.NewT100(
		req.DstAppID,
		req.SrcAppID,
		req.AppClientVersion,
		req.MainSigMap,
		h.opt.Client.SSOVersion,
	)
	tlvs[0x010a] = tlv.NewT10A(sig.Tickets["A2"].Sig)
	tlvs[0x0116] = tlv.NewT116(
		h.opt.Client.MiscBitmap,
		req.SubSigMap,
		req.SubAppIDList,
	)
	tlvs[0x0108] = tlv.NewT108(sig.Session.KSID)
	key := [16]byte{}
	if !req.changeD2 {
		copy(key[:], sig.Tickets["A2"].Key)
	} else {
		key = md5.Sum(sig.Tickets["D2"].Key)
	}
	tlvs[0x0144] = tlv.NewT144(key,
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
	if !req.changeD2 {
		tlvs[0x0145] = tlv.NewT145(util.BytesToSTBytes(h.opt.Device.GUID))
	} else {
		tlvs[0x0143] = tlv.NewT143(sig.Tickets["D2"].Sig)
	}
	tlvs[0x0142] = tlv.NewT142([]byte(h.opt.Client.PackageName))
	tlvs[0x0154] = tlv.NewT154(req.GetSeq())
	tlvs[0x0018] = tlv.NewT18(
		req.DstAppID,
		req.MainSigMap,
		req.GetUin(),
		0x0000,
	)
	tlvs[0x0141] = tlv.NewT141(
		defaultDeviceSIMOPName,
		h.opt.Device.NetworkType,
		defaultDeviceAPNName,
	)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	if len(req.Domains) > 0 {
		tlvs[0x0511] = tlv.NewT511(req.Domains)
	}
	tlvs[0x0147] = tlv.NewT147(
		req.DstAppID,
		[]byte(h.opt.Client.VersionName),
		util.BytesToSTBytes(h.opt.Client.SignatureMD5),
	)
	if len(h.t172) != 0 {
		tlvs[0x0172] = tlv.NewT172(h.t172)
	}
	tlvs[0x0177] = tlv.NewT177(h.opt.Client.BuildTime, h.opt.Client.SDKVersion)
	tlvs[0x0187] = tlv.NewT187(md5.Sum([]byte(h.opt.Device.MACAddress)))
	tlvs[0x0188] = tlv.NewT188(md5.Sum([]byte(h.opt.Device.OSBuildID)))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(h.opt.Device.IMSI)))
	// // DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(
		md5.Sum([]byte(h.opt.Device.BSSIDAddress)),
		[]byte(h.opt.Device.SSIDAddress),
	)
	// // DISABLED: tgt
	// tlvs[0x0544] = tlv.NewT544(
	// 	req.Uin,
	// 	h.opt.Device.GUID,
	// 	h.opt.Client.SDKVersion,
	// 	0x0009,
	// )
	if !req.changeD2 {
		req.SetType(0x000a)
	} else {
		req.SetType(0x000b)
	}
	req.SetServiceMethod(ServiceMethodAuthExchangeAccount)
	return tlvs
}

func (h *Handler) getSessionTicketsWithoutPassword(
	ctx context.Context,
	req *getSessionTicketsWithoutPasswordRequest,
) (*Response, error) {
	return h.getSessionTickets(ctx, req)
}
