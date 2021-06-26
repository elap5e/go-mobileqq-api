package oicq

import (
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type Message struct {
	Version       uint16
	ServiceMethod uint16
	Uin           uint64
	EncryptMethod uint8
	RandomKey     [16]byte
	KeyVersion    uint16
	PublicKey     []byte
	ShareKey      [16]byte
	Type          uint16
	Code          uint8
	TLVs          map[uint16]tlv.TLVCodec
}
