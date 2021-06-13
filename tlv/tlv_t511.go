package tlv

import (
	"log"
	"strconv"
	"strings"

	"github.com/elap5e/go-mobileqq-api/bytes"
)

type T511 struct {
	tlv  *TLV
	list []string
}

func NewT511(list []string) *T511 {
	return &T511{
		tlv:  NewTLV(0x0511, 0x0000, nil),
		list: list,
	}
}

func (t *T511) Encode(b *bytes.Buffer) {
	v := bytes.NewBuffer([]byte{})
	var list []string
	for i := range t.list {
		if t.list[i] != "" {
			list = append(list, t.list[i])
		}
	}
	v.EncodeUint16(uint16(len(list)))
	var flag uint8
	for _, str := range list {
		idx0 := strings.Index(str, "(")
		idx1 := strings.Index(str, ")")
		if idx0 != 0 || idx1 <= 0 {
			flag = 0x01
		} else {
			i, err := strconv.Atoi(str[idx0+1 : idx1])
			if err != nil {
				log.Printf("GetTLV0x0511 error: %s", err.Error())
			}
			var z1 = (1048576 & i) > 0
			var z2 = (i & 134217728) > 0
			if z1 {
				flag = 0x01
			} else {
				flag = 0x00
			}
			if z2 {
				flag |= 0x02
			}
			str = str[idx1+1:]
		}
		v.EncodeUint8(flag)
		v.EncodeString(str)
	}
	t.tlv.SetValue(v)
	t.tlv.Encode(b)
}

func (t *T511) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	_, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}
