package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

var (
	deviceBootloader   = []byte("unknown")
	deviceProcVersion  = []byte("Linux version 2.6.18-92.el5 (brewbuilder@ls20-bc2-13.build.redhat.com)")
	deviceCodename     = []byte("davinci")
	deviceIncremental  = []byte("20.10.20")
	deviceFingerprint  = []byte("Xiaomi/davinci/davinci:11/RKQ1.200827.002/20.10.20:user/release-keys")
	deviceBootID       = []byte("aa6bf49c-a995-4761-874f-8b1a9eee341e")
	deviceOSBuildID    = []byte("RKQ1.200827.002")
	deviceBaseband     = []byte("4.3.c5-00069-SM6150_GEN_PACK-1")
	deviceInnerVersion = []byte("20.10.20")

	ssoVersion = uint32(0x00000011)
)

func SetSSOVersion(ver uint32) {
	ssoVersion = ver
}

func SetDeviceOSBuildID(id []byte) {
	deviceOSBuildID = id
}

type TLV struct {
	t uint16
	l uint16
	v *bytes.Buffer
}

type TLVCodec interface {
	Encode(b *bytes.Buffer)
	Decode(b *bytes.Buffer) error
}

func NewTLV(t uint16, l uint16, v *bytes.Buffer) *TLV {
	return &TLV{t: t, l: l, v: v}
}

func (t *TLV) SetValue(v *bytes.Buffer) {
	t.v = v
}

func (t *TLV) GetType() uint16 {
	return t.t
}

func (t *TLV) GetValue() (*bytes.Buffer, error) {
	return t.v, nil
}

func (t *TLV) MustGetValue() *bytes.Buffer {
	return t.v
}

func (t *TLV) Encode(b *bytes.Buffer) {
	v := t.v.Bytes()
	t.l = uint16(len(v))
	b.EncodeUint16(t.t)
	b.EncodeUint16(t.l)
	b.EncodeRawBytes(v)
}

func (t *TLV) Decode(b *bytes.Buffer) error {
	var err error
	if t.t, err = b.DecodeUint16(); err != nil {
		return err
	}
	var v []byte
	if v, err = b.DecodeBytes(); err != nil {
		return err
	}
	t.l = uint16(len(v))
	t.v = bytes.NewBuffer(v)
	return nil
}
