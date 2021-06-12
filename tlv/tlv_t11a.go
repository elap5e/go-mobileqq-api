package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T11A struct {
	tlv    *TLV
	face   []byte
	age    []byte
	gender []byte
	nick   []byte
}

func NewT11A(face, age, gender, nick []byte) *T11A {
	return &T11A{
		tlv:    NewTLV(0x011a, 0x0000, nil),
		face:   face,
		age:    age,
		gender: gender,
		nick:   nick,
	}
}

func (t *T11A) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeRawBytes(t.face)
	v.EncodeRawBytes(t.age)
	v.EncodeRawBytes(t.gender)
	v.EncodeBytes(t.nick)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T11A) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.face, err = v.DecodeBytesN(0x0002); err != nil {
		return err
	}
	if t.age, err = v.DecodeBytesN(0x0001); err != nil {
		return err
	}
	if t.gender, err = v.DecodeBytesN(0x0001); err != nil {
		return err
	}
	if t.nick, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T11A) GetFace() ([]byte, error) {
	return t.face, nil
}

func (t *T11A) GetAge() ([]byte, error) {
	return t.age, nil
}

func (t *T11A) GetGender() ([]byte, error) {
	return t.gender, nil
}

func (t *T11A) GetNick() ([]byte, error) {
	return t.nick, nil
}
