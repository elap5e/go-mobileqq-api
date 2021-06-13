package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T149 struct {
	tlv         *TLV
	contentType uint16
	title       []byte
	content     []byte
	otherInfo   []byte
}

func NewT149(contentType uint16, title, content, otherInfo []byte) *T149 {
	return &T149{
		tlv:         NewTLV(0x0149, 0x0000, nil),
		contentType: contentType,
		title:       title,
		content:     content,
		otherInfo:   otherInfo,
	}
}

func (t *T149) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(t.contentType)
	v.EncodeBytes(t.title)
	v.EncodeBytes(t.content)
	v.EncodeBytes(t.otherInfo)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T149) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.contentType, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.title, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.content, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.otherInfo, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T149) GetContentType() (uint16, error) {
	return t.contentType, nil
}

func (t *T149) GetTitle() ([]byte, error) {
	return t.title, nil
}

func (t *T149) GetContent() ([]byte, error) {
	return t.content, nil
}

func (t *T149) GetOtherInfo() ([]byte, error) {
	return t.otherInfo, nil
}
