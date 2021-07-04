package rpc

import (
	"context"
	"crypto/md5"
	"log"
	"path"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthGetSessionTicketsRequest interface {
	GetTLVs(ctx context.Context) (map[uint16]tlv.TLVCodec, error)
	GetType() uint16
	SetType(typ uint16)
	GetSeq() uint32
	SetSeq(seq uint32)
	GetUin() uint64
	GetUsername() string
	SetUsername(username string)
	GetServiceMethod() string
}

type authGetSessionTicketsRequest struct {
	typ      uint16
	seq      uint32
	uin      uint64
	username string
	method   string
}

func (req *authGetSessionTicketsRequest) GetType() uint16 {
	return req.typ
}

func (req *authGetSessionTicketsRequest) SetType(typ uint16) {
	req.typ = typ
}

func (req *authGetSessionTicketsRequest) GetSeq() uint32 {
	return req.seq
}

func (req *authGetSessionTicketsRequest) SetSeq(seq uint32) {
	req.seq = seq
}

func (req *authGetSessionTicketsRequest) GetUin() uint64 {
	return req.uin
}

func (req *authGetSessionTicketsRequest) GetUsername() string {
	return req.username
}

func (req *authGetSessionTicketsRequest) SetUsername(username string) {
	req.username = username
	uin, _ := strconv.ParseInt(username, 10, 64)
	req.uin = uint64(uin)
}

func (req *authGetSessionTicketsRequest) GetServiceMethod() string {
	return req.method
}

func (req *authGetSessionTicketsRequest) SetServiceMethod(method string) {
	req.method = method
}

type AuthGetSessionTicketsResponse struct {
	Code     uint8
	Uin      uint64
	Username string

	PictureSign  []byte
	PictureData  []byte
	CaptchaSign  string
	ErrorCode    uint32
	ErrorTitle   string
	ErrorMessage string
	Message      string
	SMSMobile    string

	AuthSession []byte
	T119        []byte
	T150        []byte
	T161        []byte
	T174        []byte
	T17B        []byte
	T401        [16]byte
	T402        []byte
	T403        []byte
	T546        []byte

	LoginExtraData []byte
}

func (resp *AuthGetSessionTicketsResponse) SetTLVs(
	ctx context.Context,
	tlvs map[uint16]tlv.TLVCodec,
) error {
	if v, ok := tlvs[0x0104].(*tlv.TLV); ok {
		resp.AuthSession = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0105].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		resp.PictureSign, _ = buf.DecodeBytes()
		resp.PictureData, _ = buf.DecodeBytes()
	}
	if v, ok := tlvs[0x0119].(*tlv.TLV); ok {
		resp.T119 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0174].(*tlv.TLV); ok {
		resp.T174 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x017b].(*tlv.TLV); ok {
		resp.T17B = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0192].(*tlv.TLV); ok {
		resp.CaptchaSign = string(v.MustGetValue().Bytes())
	}
	if v, ok := tlvs[0x0146].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		resp.ErrorCode, _ = buf.DecodeUint32()
		resp.ErrorTitle, _ = buf.DecodeString()
		resp.ErrorMessage, _ = buf.DecodeString()
	}
	if v, ok := tlvs[0x0150].(*tlv.TLV); ok {
		resp.T150 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0161].(*tlv.TLV); ok {
		resp.T161 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x017e].(*tlv.TLV); ok {
		resp.Message = string(v.MustGetValue().Bytes())
	}
	if v, ok := tlvs[0x0178].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		_, _ = buf.DecodeString()
		mobile, _ := buf.DecodeString()
		resp.SMSMobile = mobile
	}
	if v, ok := tlvs[0x0402].(*tlv.TLV); ok {
		resp.T402 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0403].(*tlv.TLV); ok {
		resp.T403 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0537].(*tlv.TLV); ok {
		resp.LoginExtraData, _ = v.MustGetValue().DecodeBytes()
	}
	if v, ok := tlvs[0x0546].(*tlv.TLV); ok {
		resp.T546 = v.MustGetValue().Bytes()
	}
	return nil
}

func (c *Client) AuthGetSessionTickets(
	ctx context.Context,
	req AuthGetSessionTicketsRequest,
) (*AuthGetSessionTicketsResponse, error) {
	req.SetSeq(c.getNextSeq())
	tlvs, err := req.GetTLVs(ctx)
	if err != nil {
		return nil, err
	}
	buf, err := oicq.Marshal(ctx, &oicq.Message{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.GetUin(),
		EncryptMethod: oicq.EncryptMethodECDH,
		RandomKey:     c.randomKey,
		KeyVersion:    c.serverPublicKeyVersion,
		PublicKey:     c.privateKey.Public().Bytes(),
		ShareKey:      c.privateKey.ShareKey(c.serverPublicKey),
		Type:          req.GetType(),
		TLVs:          tlvs,
	})
	if err != nil {
		return nil, err
	}
	s2c := codec.ServerToClientMessage{}
	if err := c.Call(req.GetServiceMethod(), &codec.ClientToServerMessage{
		Username: req.GetUsername(),
		Seq:      req.GetSeq(),
		Buffer:   buf,
		Simple:   false,
	}, &s2c); err != nil {
		return nil, err
	}
	resp := AuthGetSessionTicketsResponse{}
	msg := &oicq.Message{
		RandomKey:  c.randomKey,
		KeyVersion: c.serverPublicKeyVersion,
		PublicKey:  c.privateKey.Public().Bytes(),
		ShareKey:   c.privateKey.ShareKey(c.serverPublicKey),
	}
	if err := oicq.Unmarshal(ctx, s2c.Buffer, msg); err != nil {
		return nil, err
	}
	resp.Code = msg.Code
	resp.Username = strconv.Itoa(int(msg.Uin))
	resp.Uin = msg.Uin
	if err := resp.SetTLVs(ctx, msg.TLVs); err != nil {
		return nil, err
	}
	switch resp.Code {
	case 0x00:
		// success
		c.loginExtraData = resp.LoginExtraData

		// decode t119
		sig := c.GetUserSignature(req.GetUsername())
		key := [16]byte{}
		switch msg.Type {
		default:
			log.Printf("x_x [oicq] type:0x%04x", msg.Type)
			copy(key[:], sig.Tickets["A1"].Key)
		case 0x0009, 0x000f, 0x0014: // ???
			copy(key[:], sig.Tickets["A1"].Key)
		case 0x000a:
			copy(key[:], sig.Tickets["A2"].Key)
		case 0x000b:
			key = md5.Sum(sig.Tickets["D2"].Key)
		}
		t119, err := crypto.NewCipher(key).Decrypt(resp.T119)
		if err != nil {
			return nil, err
		}

		// log.Printf("--> [recv] dump tlv 0x0119(decrypt):\n%s", hex.Dump(t119))
		tlvs := map[uint16]tlv.TLVCodec{}
		buf := bytes.NewBuffer(t119)
		l, _ := buf.DecodeUint16()
		for i := 0; i < int(l); i++ {
			v := tlv.TLV{}
			v.Decode(buf)
			tlvs[v.GetType()] = &v
		}
		if c.cfg.LogLevel&LogLevelTrace != 0 {
			tlv.DumpTLVs(ctx, tlvs)
		}

		c.SetUserSignature(ctx, resp.Username, tlvs)
		c.SetUserAuthSession(resp.Username, nil)
		if v, ok := tlvs[0x0108]; ok {
			c.SetUserKSIDSession(
				resp.Username,
				v.(*tlv.TLV).MustGetValue().Bytes(),
			)
		}
		c.SaveUserSignatures(path.Join(
			c.cfg.BaseDir, PATH_TO_USER_SIGNATURE_JSON,
		))

		log.Printf(
			"^_^ [info] login success, uin:%s code:0x00",
			resp.Username,
		)
	case 0x02:
		// captcha
		c.SetUserAuthSession(resp.Username, resp.AuthSession)

		c.extraData[0x0547] = resp.T546 // TODO: check
		if resp.CaptchaSign != "" {
			log.Printf(
				">_x [warn] need captcha verify, uin %s, url %s, code 0x02",
				resp.Username, resp.CaptchaSign,
			)
		} else {
			log.Printf(
				">_x [warn] need picture verify, uin %s, code 0x02",
				resp.Username,
			)
		}
	case 0xa0:
		// device lock
		c.SetUserAuthSession(resp.Username, resp.AuthSession)

		c.t17b = resp.T17B
		log.Printf(
			">_x [warn] need sms mobile verify response, uin %s, code 0xa0",
			resp.Username,
		)
	case 0xef:
		// device lock
		c.SetUserAuthSession(resp.Username, resp.AuthSession)

		c.t174 = resp.T174
		c.t402 = resp.T402
		c.t403 = resp.T403
		c.hashedGUID = md5.Sum(
			append(append(
				c.cfg.Device.GUID,
				c.randomPassword[:]...),
				c.t402...),
		)
		log.Printf(
			">_x [warn] need sms mobile verify, uin %s, mobile %s, code 0x%02x, message %s, code 0xef",
			resp.Username, resp.SMSMobile, resp.Code, resp.Message,
		)
	case 0x01:
		log.Printf(
			"x_x [fail] invalid login, uin %s, code 0x01, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0x06:
		log.Printf(
			"x_x [fail] not implement, uin %s, code 0x06, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0x09:
		log.Printf(
			"x_x [fail] invalid service, uin %s, code 0x09, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0x10:
		log.Printf(
			"x_x [fail] session expired, uin %s, code 0x10, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0x28:
		log.Printf(
			"x_x [fail] protection mode, uin %s, code 0x10, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0xed:
		log.Printf(
			"x_x [fail] invalid device, uin %s, code 0xed, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0xcc:
		c.SetUserAuthSession(resp.Username, resp.AuthSession)

		c.t402 = resp.T402
		c.t403 = resp.T403
		c.hashedGUID = md5.Sum(
			append(append(
				c.cfg.Device.GUID,
				c.randomPassword[:]...),
				c.t402...),
		)
		return c.AuthUnlockDevice(ctx, NewAuthUnlockDeviceRequest(
			resp.Username,
		))
	}
	return &resp, nil
}
