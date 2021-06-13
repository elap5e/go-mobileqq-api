package tlv

import (
	"crypto/md5"
	"encoding/binary"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T184 struct {
	tlv *TLV
	j   uint64
	str string
}

func NewT184(j uint64, str string) *T184 {
	return &T184{
		tlv: NewTLV(0x0184, 0x0000, nil),
		j:   j,
		str: str,
	}
}

func (t *T184) Encode(b *bytes.Buffer) {
	v := md5.Sum([]byte(t.str))
	tmp := append(v[:], make([]byte, 8)...)
	binary.BigEndian.PutUint64(tmp[16:], t.j)
	v = md5.Sum(tmp)
	t.tlv.SetValue(bytes.NewBuffer(v[:]))
	t.tlv.Encode(b)
}

func (t *T184) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	_, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}
