package client

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/rpc"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

const PATH_TO_USER_SIGNATURE_JSON = "user_signatures.json"

func (c *Client) initUserSignatures() {
	c.userSignatures = make(map[string]*rpc.UserSignature)
	c.LoadUserSignatures(path.Join(c.cfg.BaseDir, PATH_TO_USER_SIGNATURE_JSON))

}

func (c *Client) LoadUserSignatures(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	sigs := make(map[string]*rpc.UserSignature)
	if err = json.Unmarshal(data, &sigs); err != nil {
		return err
	}
	for username, sig := range sigs {
		tsig := c.GetUserSignature(username)
		if len(sig.DeviceToken) != 0 {
			tsig.DeviceToken = sig.DeviceToken
		}
		for k, v := range sig.Domains {
			tsig.Domains[k] = v
		}
		for k, v := range sig.Tickets {
			tsig.Tickets[k] = v
		}
		tsig.Session.Auth = sig.Session.Auth
		tsig.Session.Cookie = sig.Session.Cookie
		tsig.Session.KSID = sig.Session.KSID
	}
	return nil
}

func (c *Client) SaveUserSignatures(file string) error {
	c.userSignaturesMux.Lock()
	defer c.userSignaturesMux.Unlock()
	data, err := json.MarshalIndent(c.userSignatures, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, append(data, '\n'), 0600)
}

func (c *Client) GetUserSignature(username string) *rpc.UserSignature {
	sig, ok := c.userSignatures[username]
	if !ok {
		sig = &rpc.UserSignature{}
		sig.Username = username
		sig.Domains = make(map[string]string)
		sig.Tickets = make(map[string]*rpc.Ticket)
		sig.Session = &rpc.Session{}
		sig.Tickets["A1"] = &rpc.Ticket{
			Sig: nil,
			Key: make([]byte, 16),
			Iss: time.Now().Unix(),
			Exp: 0,
		}
		rand.Read(sig.Tickets["A1"].Key)
		sig.Session.Cookie = make([]byte, 4)
		rand.Read(sig.Session.Cookie)
		c.userSignatures[username] = sig
		c.rpc.SetUserSignature(username, sig)
		cacheDir := c.GetCacheByUsernameDir(username)
		_, err := os.Stat(cacheDir)
		if os.IsNotExist(err) {
			err = os.Mkdir(cacheDir, 0755)
		}
		if err != nil {
			log.Fatal().Err(err).
				Msg("failed to mkdir " + cacheDir)
		}
	}
	if c.db != nil {
		uin, _ := strconv.ParseUint(username, 10, 64)
		if err := c.dbCreateChannelTableByUin(uin); err != nil {
			log.Fatal().Err(err).
				Msg("failed to operate database")
		}
		if err := c.dbCreateContactTableByUin(uin); err != nil {
			log.Fatal().Err(err).
				Msg("failed to operate database")
		}
		if err := c.dbCreateMessageRecordTableByUin(uin); err != nil {
			log.Fatal().Err(err).
				Msg("failed to operate database")
		}
	}
	return sig
}

func (c *Client) SetUserSignature(
	ctx context.Context,
	username string,
	tlvs map[uint16]tlv.TLVCodec,
) {
	sig := c.GetUserSignature(username)
	tsig := ParseUserSignature(ctx, username, tlvs)
	if len(tsig.DeviceToken) != 0 {
		sig.DeviceToken = tsig.DeviceToken
	}
	for key, val := range tsig.Domains {
		sig.Domains[key] = val
	}
	for key, val := range tsig.Tickets {
		if ssig, ok := sig.Tickets[key]; ok {
			if len(val.Sig) != 0 {
				ssig.Sig = val.Sig
			}
			if len(val.Key) != 0 {
				ssig.Key = val.Key
			}
			if val.Iss != 0 {
				ssig.Iss = val.Iss
			}
			if val.Exp != 0 {
				ssig.Exp = val.Exp
			}
		} else {
			sig.Tickets[key] = val
		}
	}
}

func (c *Client) SetUserAuthSession(username string, session []byte) {
	sig := c.GetUserSignature(username)
	sig.Session.Auth = session
}

func (c *Client) SetUserKSIDSession(username string, ksid []byte) {
	sig := c.GetUserSignature(username)
	sig.Session.KSID = ksid
}

func ParseUserSignature(
	ctx context.Context,
	username string,
	tlvs map[uint16]tlv.TLVCodec,
) *rpc.UserSignature {
	token := []byte{}
	if v, ok := tlvs[0x0322]; ok {
		token = v.(*tlv.TLV).MustGetValue().Bytes()
	}

	domains := map[string]string{}
	if v, ok := tlvs[0x0512]; ok {
		buf := bytes.NewBuffer(v.(*tlv.TLV).MustGetValue().Bytes())
		l, _ := buf.DecodeUint16()
		for i := 0; i < int(l); i++ {
			key, _ := buf.DecodeString()
			domains[key], _ = buf.DecodeString()
			_, _ = buf.DecodeUint16()
		}
	}

	chgt := map[uint16]uint32{}
	if v, ok := tlvs[0x0138]; ok {
		buf := bytes.NewBuffer(v.(*tlv.TLV).MustGetValue().Bytes())
		l, _ := buf.DecodeUint32()
		for i := 0; i < int(l); i++ {
			key, _ := buf.DecodeUint16()
			chgt[key], _ = buf.DecodeUint32()
			_, _ = buf.DecodeUint32()
		}
	}

	tickets := map[string]*rpc.Ticket{}
	{
		if v, ok := tlvs[0x0106]; ok {
			tickets["A1"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: tlvs[0x010c].(*tlv.TLV).MustGetValue().Bytes(),
				Iss: time.Now().Unix(),
				Exp: 0,
			}
		}
		if v, ok := tlvs[0x010a]; ok {
			tickets["A2"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: tlvs[0x010d].(*tlv.TLV).MustGetValue().Bytes(),
				Iss: time.Now().Unix(),
				Exp: int64(chgt[0x010a]),
			}
		}
		if v, ok := tlvs[0x010b]; ok {
			tickets["A5"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: 0,
			}
		}
		if v, ok := tlvs[0x0102]; ok {
			tickets["A8"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: 0,
			}
		}
		if v, ok := tlvs[0x0143]; ok {
			tickets["D2"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: tlvs[0x0305].(*tlv.TLV).MustGetValue().Bytes(),
				Iss: time.Now().Unix(),
				Exp: int64(chgt[0x0143]),
			}
		}
		if v, ok := tlvs[0x011c]; ok {
			tickets["LSKey"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: int64(chgt[0x011c]),
			}
		}
		if v, ok := tlvs[0x0120]; ok {
			tickets["SKey"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: int64(chgt[0x0120]),
			}
		}
		if v, ok := tlvs[0x0121]; ok {
			tickets["Sig64"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: 0,
			}
		}
		if v, ok := tlvs[0x0164]; ok {
			tickets["SID"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: int64(chgt[0x0164]),
			}
		}
		if v, ok := tlvs[0x0114]; ok {
			tickets["ST"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: tlvs[0x010e].(*tlv.TLV).MustGetValue().Bytes(),
				Iss: time.Now().Unix(),
				Exp: 0,
			}
		}
		if v, ok := tlvs[0x0103]; ok {
			tickets["STWeb"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: int64(chgt[0x0103]),
			}
		}
		if v, ok := tlvs[0x016d]; ok {
			tickets["SuperKey"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: 0,
			}
		}
		if v, ok := tlvs[0x0136]; ok {
			tickets["VKey"] = &rpc.Ticket{
				Sig: v.(*tlv.TLV).MustGetValue().Bytes(),
				Key: nil,
				Iss: time.Now().Unix(),
				Exp: int64(chgt[0x0136]),
			}
		}
		// 0x00004000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: 0,
		// }, // ??? OpenKey 0x0125
		// 0x00008000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: 0,
		// }, // ??? AccessToken 0x0132
		// 0x00200000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: 0,
		// }, // ??? AQSig 0x0171
		// 0x00800000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: 0,
		// }, // ??? PayToken 0x0199
		// 0x01000000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: 0,
		// }, // ??? PF 0x0200
		// 0x02000000: {
		// 	Sig: nil,
		// 	Key: nil,
		// 	Iss: time.Now().Unix(),
		// 	Exp: 0,
		// }, // ??? DA2 0x0203
	}

	return &rpc.UserSignature{
		Username:    username,
		DeviceToken: token,
		Domains:     domains,
		Tickets:     tickets,
	}
}
