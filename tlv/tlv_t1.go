package tlv

import (
	"fmt"
	"math/rand"
	"net"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/util"
)

type T1 struct {
	tlv *TLV
	uin uint64
	ip  net.IP
}

func NewT1(uin uint64, ip net.IP) *T1 {
	return &T1{
		tlv: NewTLV(0x0001, 0x0014, nil),
		uin: uin,
		ip:  ip,
	}
}

func (t *T1) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(0x0001)
	v.EncodeUint32(rand.Uint32())
	v.EncodeUint32(uint32(t.uin))
	v.EncodeUint32(util.GetServerTime())
	v.EncodeRawBytes(t.ip.To4())
	v.EncodeUint16(0x0000)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T1) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	ipVer, err := v.DecodeUint16()
	if err != nil {
		return err
	} else if ipVer != 0x0001 {
		return fmt.Errorf("ip version 0x%x not support", ipVer)
	}
	if _, err := v.DecodeUint32(); err != nil {
		return err
	}
	uin, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.uin = uint64(uin)
	if _, err := v.DecodeUint32(); err != nil {
		return err
	}
	if t.ip, err = v.DecodeBytesN(0x0004); err != nil {
		return err
	}
	if _, err := v.DecodeUint32(); err != nil {
		return err
	}
	return nil
}

func (t *T1) GetUin() (uint64, error) {
	return t.uin, nil
}

func (t *T1) GetIP() (net.IP, error) {
	return t.ip, nil
}
