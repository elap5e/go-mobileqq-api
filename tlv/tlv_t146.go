package tlv

import (
	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T146 struct {
	tlv         *TLV
	version     uint16
	code        uint16
	title       []byte
	message     []byte
	messageType uint16
	errorInfo   []byte
}

func NewT146(version uint16, code uint16, title []byte, message []byte, messageType uint16, errorInfo []byte) *T146 {
	return &T146{
		tlv:         NewTLV(0x0146, 0x0000, nil),
		version:     version,
		code:        code,
		title:       title,
		message:     message,
		messageType: messageType,
		errorInfo:   errorInfo,
	}
}

func (t *T146) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	v.EncodeUint16(t.version)
	v.EncodeUint16(t.code)
	v.EncodeBytes(t.title)
	v.EncodeBytes(t.message)
	v.EncodeUint16(t.messageType)
	v.EncodeBytes(t.errorInfo)
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T146) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.version, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.code, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.title, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.message, err = v.DecodeBytes(); err != nil {
		return err
	}
	if t.messageType, err = v.DecodeUint16(); err != nil {
		return err
	}
	if t.errorInfo, err = v.DecodeBytes(); err != nil {
		return err
	}
	return nil
}

func (t *T146) GetVersion() (uint16, error) {
	return t.version, nil
}

func (t *T146) GetCode() (uint16, error) {
	return t.code, nil
}

func (t *T146) GetTitle() ([]byte, error) {
	return t.title, nil
}

func (t *T146) GetMessage() ([]byte, error) {
	return t.message, nil
}

func (t *T146) GetMessageType() (uint16, error) {
	return t.messageType, nil
}

func (t *T146) GetErrorInfo() ([]byte, error) {
	return t.errorInfo, nil
}
