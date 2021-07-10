package client

import (
	"context"
	"crypto/md5"
	"fmt"
	"path"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/tlv"
	"github.com/rs/zerolog"
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
	resp.Username = strconv.FormatUint(msg.Uin, 10)
	resp.Uin = msg.Uin
	if err := resp.SetTLVs(ctx, msg.TLVs); err != nil {
		return nil, err
	}
	switch resp.Code {
	default:
		log.Warn().
			Str("@uin", resp.Username).
			Uint8("code", resp.Code).
			Msg("--> [auth] not implement")
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
		case 0x0007: // AuthCheckSMSAndGetSessionTickets
			copy(key[:], sig.Tickets["A1"].Key)
		case 0x0009: // AuthGetSessionTicketsWithPassword
			copy(key[:], sig.Tickets["A1"].Key)
		case 0x000a: // AuthGetSessionTicketsWithoutPassword.A2
			copy(key[:], sig.Tickets["A2"].Key)
		case 0x000b: // AuthGetSessionTicketsWithoutPassword.D2
			key = md5.Sum(sig.Tickets["D2"].Key)
		case 0x0014: // AuthUnlockDevice
			copy(key[:], sig.Tickets["A1"].Key)
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
		if log.GetLevel() == zerolog.TraceLevel {
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

		log.Info().
			Str("@uin", resp.Username).
			Uint8("code", resp.Code).
			Msg("^.^ [auth] login success")
	case 0x02:
		// captcha
		c.SetUserAuthSession(resp.Username, resp.AuthSession)

		c.extraData[0x0547] = resp.T546 // TODO: check
		if resp.CaptchaSign != "" {
			log.Warn().
				Str("@uin", resp.Username).
				Uint8("code", resp.Code).
				Msg("x<- [auth] need captcha verify, url " + resp.CaptchaSign)
		} else {
			log.Warn().
				Str("@uin", resp.Username).
				Uint8("code", resp.Code).
				Msg("x<- [auth] need picture verify")
		}
	case 0xa0:
		// device lock
		c.SetUserAuthSession(resp.Username, resp.AuthSession)

		c.t17b = resp.T17B
		log.Warn().
			Str("@uin", resp.Username).
			Uint8("code", resp.Code).
			Msg("x<- [auth] need sms verify response")
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
		log.Warn().
			Str("@uin", resp.Username).
			Uint8("code", resp.Code).
			Str("info", resp.Message).
			Msg("x<- [auth] need sms verify, mobile " + resp.SMSMobile)
	case 0x01:
		log.Error().
			Err(fmt.Errorf("%s: %s", resp.ErrorTitle, resp.ErrorMessage)).
			Str("@uin", resp.Username).
			Uint8("code", resp.Code).
			Msg("x-x [auth] invalid login")
	case 0x06:
		log.Error().
			Err(fmt.Errorf("%s: %s", resp.ErrorTitle, resp.ErrorMessage)).
			Str("@uin", resp.Username).
			Uint8("code", resp.Code).
			Msg("x-x [auth] not implement")
	case 0x09:
		log.Error().
			Err(fmt.Errorf("%s: %s", resp.ErrorTitle, resp.ErrorMessage)).
			Str("@uin", resp.Username).
			Uint8("code", resp.Code).
			Msg("x-x [auth] invalid service")
	case 0x10:
		log.Error().
			Err(fmt.Errorf("%s: %s", resp.ErrorTitle, resp.ErrorMessage)).
			Str("@uin", resp.Username).
			Uint8("code", resp.Code).
			Msg("x-x [auth] session expired")
	case 0x28:
		log.Error().
			Err(fmt.Errorf("%s: %s", resp.ErrorTitle, resp.ErrorMessage)).
			Str("@uin", resp.Username).
			Uint8("code", resp.Code).
			Msg("x-x [auth] protection mode")
	case 0xed:
		log.Error().
			Err(fmt.Errorf("%s: %s", resp.ErrorTitle, resp.ErrorMessage)).
			Str("@uin", resp.Username).
			Uint8("code", resp.Code).
			Msg("x-x [auth] invalid device")
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
