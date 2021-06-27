package oicq

import (
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type Message struct {
	Version       uint16
	ServiceMethod uint16
	Uin           uint64
	EncryptMethod EncryptMethod
	RandomKey     [16]byte
	KeyVersion    uint16
	PublicKey     []byte
	ShareKey      [16]byte
	Type          uint16
	Code          uint8
	TLVs          map[uint16]tlv.TLVCodec
}

type EncryptMethod uint8

var (
	EncryptMethod0x00 EncryptMethod = 0x00
	EncryptMethod0x03 EncryptMethod = 0x03
	EncryptMethodECDH EncryptMethod = 0x07 | 0x80 // 0x07: no password login?
	EncryptMethodST   EncryptMethod = 0x45
	EncryptMethodNULL EncryptMethod = 0xff
)

func GetEncryptMethod(v uint8) EncryptMethod {
	switch v {
	case 0x00:
		return EncryptMethod0x00
	case 0x03:
		return EncryptMethod0x00
	case 0x07, 0x87:
		return EncryptMethodECDH
	case 0x45:
		return EncryptMethodST
	default:
		return EncryptMethodNULL
	}
}
