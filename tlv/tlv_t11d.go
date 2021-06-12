package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T11D struct {
	tlv   *TLV
	appID uint64
	stKey []byte
	st    []byte
}

func NewT11D(appID uint64, stKey, st []byte) *T11D {
	return &T11D{
		tlv:   NewTLV(0x011d, 0x0000, nil),
		appID: appID,
		stKey: stKey,
		st:    st,
	}
}

func (t *T11D) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint32(uint32(t.appID))
	v.EncodeRawBytes(t.stKey)
	v.EncodeBytes(t.st)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T11D) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	appID, err := v.DecodeUint32()
	if err != nil {
		return err
	}
	t.appID = uint64(appID)
	if t.stKey, err = v.DecodeBytesN(0x0010); err != nil {
		return err
	}
	if t.st, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T11D) GetAppID() (uint64, error) {
	return t.appID, nil
}

func (t *T11D) GetSTKey() ([]byte, error) {
	return t.stKey, nil
}

func (t *T11D) GetST() ([]byte, error) {
	return t.st, nil
}
