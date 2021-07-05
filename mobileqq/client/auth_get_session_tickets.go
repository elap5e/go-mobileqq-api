package client

import (
	"context"
	"crypto/md5"
	"path"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

func (c *Client) AuthGetSessionTickets(
	ctx context.Context,
	req AuthGetSessionTicketsRequest,
) (*AuthGetSessionTicketsResponse, error) {
	req.SetSeq(c.rpc.GetNextSeq())
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
	if err := c.rpc.Call(req.GetServiceMethod(), &codec.ClientToServerMessage{
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
			log.Error().Msgf("x_x [oicq] type:0x%04x", msg.Type)
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

		tlvs := map[uint16]tlv.TLVCodec{}
		buf := bytes.NewBuffer(t119)
		l, _ := buf.DecodeUint16()
		for i := 0; i < int(l); i++ {
			v := tlv.TLV{}
			v.Decode(buf)
			tlvs[v.GetType()] = &v
		}
		// tlv.DumpTLVs(ctx, tlvs)

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

		log.Info().Msgf(
			"^_^ [auth] login success, uin:%s code:0x00",
			resp.Username,
		)
	case 0x02:
		// captcha
		c.SetUserAuthSession(resp.Username, resp.AuthSession)

		c.extraData[0x0547] = resp.T546 // TODO: check
		if resp.CaptchaSign != "" {
			log.Warn().Msgf(
				">_x [auth] need captcha verify, uin %s, url %s, code 0x02",
				resp.Username, resp.CaptchaSign,
			)
		} else {
			log.Warn().Msgf(
				">_x [auth] need picture verify, uin %s, code 0x02",
				resp.Username,
			)
		}
	case 0xa0:
		// device lock
		c.SetUserAuthSession(resp.Username, resp.AuthSession)

		c.t17b = resp.T17B
		log.Warn().Msgf(
			">_x [auth] need sms mobile verify response, uin %s, code 0xa0",
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
		log.Warn().Msgf(
			">_x [auth] need sms mobile verify, uin %s, mobile %s, code 0x%02x, message %s, code 0xef",
			resp.Username, resp.SMSMobile, resp.Code, resp.Message,
		)
	case 0x01:
		log.Error().Msgf(
			"x_x [auth] invalid login, uin %s, code 0x01, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0x06:
		log.Error().Msgf(
			"x_x [auth] not implement, uin %s, code 0x06, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0x09:
		log.Error().Msgf(
			"x_x [auth] invalid service, uin %s, code 0x09, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0x10:
		log.Error().Msgf(
			"x_x [auth] session expired, uin %s, code 0x10, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0x28:
		log.Error().Msgf(
			"x_x [auth] protection mode, uin %s, code 0x10, error %s: %s",
			resp.Username, resp.ErrorTitle, resp.ErrorMessage,
		)
	case 0xed:
		log.Error().Msgf(
			"x_x [auth] invalid device, uin %s, code 0xed, error %s: %s",
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
