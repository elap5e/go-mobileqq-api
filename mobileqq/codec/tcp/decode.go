package tcp

import (
	_bytes "bytes"
	"compress/zlib"
	"fmt"
	"io"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (c *clientCodec) readBody(msg *codec.ServerToClientMessage) error {
	var err error
	var n uint32
	if n, err = c.bufResp.DecodeUint32(); err != nil {
		return err
	} // 0x00000000
	if msg.Seq, err = c.bufResp.DecodeUint32(); err != nil {
		return err
	}
	if msg.Code, err = c.bufResp.DecodeUint32(); err != nil {
		return err
	}
	if msg.Message, err = c.bufResp.ReadUint32String(); err != nil {
		return err
	}
	if msg.ServiceMethod, err = c.bufResp.ReadUint32String(); err != nil {
		return err
	}
	if msg.Cookie, err = c.bufResp.ReadUint32Bytes(); err != nil {
		return err
	}
	if msg.Flag, err = c.bufResp.DecodeUint32(); err != nil {
		return err
	}
	if msg.Flag != codec.FlagNoCompression &&
		msg.Flag != codec.FlagZlibCompression {
		return fmt.Errorf("unsupported flag 0x%x", msg.Flag)
	}
	if c.bufResp.Index() < int(n) {
		if msg.ReserveField, err = c.bufResp.ReadUint32Bytes(); err != nil {
			return err
		}
	}
	if msg.Buffer, err = c.bufResp.ReadUint32Bytes(); err != nil {
		return err
	}
	return nil
}

func (c *clientCodec) readHead(msg *codec.ServerToClientMessage) error {
	var err error
	if _, err = c.bufResp.ReadUint32(); err != nil {
		return err
	} // 0x00000000
	if msg.Version, err = c.bufResp.ReadUint32(); err != nil {
		return err
	}
	if msg.Version != codec.VersionDefault &&
		msg.Version != codec.VersionSimple {
		return fmt.Errorf("unsupported version 0x%x", msg.Version)
	}
	if msg.EncryptType, err = c.bufResp.DecodeUint8(); err != nil {
		return err
	}
	if msg.EncryptType != codec.EncryptTypeNotNeedEncrypt &&
		msg.EncryptType != codec.EncryptTypeEncryptByD2Key &&
		msg.EncryptType != codec.EncryptTypeEncryptByZeros {
		return fmt.Errorf("unsupported encrypt type 0x%x", msg.EncryptType)
	}
	if _, err = c.bufResp.DecodeUint8(); err != nil {
		return err
	} // 0x00
	if msg.Username, err = c.bufResp.ReadUint32String(); err != nil {
		return err
	}
	return nil
}

func (c *clientCodec) ReadBody(msg *codec.ServerToClientMessage) error {
	switch msg.EncryptType {
	case codec.EncryptTypeNotNeedEncrypt:
	case codec.EncryptTypeEncryptByD2Key:
		buf, err := crypto.NewCipher(msg.UserD2Key).Decrypt(c.bufResp.Bytes())
		if err != nil {
			return err
		}
		c.bufResp = bytes.NewBuffer(buf)
	case codec.EncryptTypeEncryptByZeros:
		buf, err := crypto.NewCipher([16]byte{}).Decrypt(c.bufResp.Bytes())
		if err != nil {
			return err
		}
		c.bufResp = bytes.NewBuffer(buf)
	}
	if err := c.readBody(msg); err != nil {
		return err
	}
	switch msg.Flag {
	case codec.FlagNoCompression:
	case codec.FlagZlibCompression:
		reader, err := zlib.NewReader(_bytes.NewReader(msg.Buffer))
		if err != nil {
			return err
		}
		defer reader.Close()
		var buf _bytes.Buffer
		io.Copy(&buf, reader)
		msg.Buffer = buf.Bytes()
	}
	return nil
}

func (c *clientCodec) ReadHead(msg *codec.ServerToClientMessage) error {
	p, err := c.read()
	if err != nil {
		return err
	}
	c.bufResp = bytes.NewBuffer(p)
	return c.readHead(msg)
}
