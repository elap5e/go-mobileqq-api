package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T167 struct {
	tlv         *TLV
	imageType   uint8
	imageFormat uint8
	imageURL    []byte
}

func NewT167(imageType, imageFormat uint8, imageURL []byte) *T167 {
	return &T167{
		tlv:         NewTLV(0x0167, 0x0000, nil),
		imageType:   imageType,
		imageFormat: imageFormat,
		imageURL:    imageURL,
	}
}

func (t *T167) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint8(t.imageType)
	v.EncodeUint8(t.imageFormat)
	v.EncodeBytes(t.imageURL)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T167) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.imageType, err = v.DecodeUint8(); err != nil {
		return err
	}
	if t.imageFormat, err = v.DecodeUint8(); err != nil {
		return err
	}
	if t.imageURL, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T167) GetImageType() (uint8, error) {
	return t.imageType, nil
}

func (t *T167) GetImageFormat() (uint8, error) {
	return t.imageFormat, nil
}

func (t *T167) GetImageURL() ([]byte, error) {
	return t.imageURL, nil
}
