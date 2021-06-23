package tlv

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type T52D struct {
	tlv *TLV
	ctx context.Context
}

func NewT52D(ctx context.Context) *T52D {
	return &T52D{
		tlv: NewTLV(0x052d, 0x0000, nil),
		ctx: ctx,
	}
}

func (t *T52D) Encode(b *bytes.Buffer) {
	v, _ := proto.Marshal(&pb.DeviceInfo{
		Bootloader:   []byte("694a7990-33fd-47c1-b550-cb30fe50f1aa"),
		Codename:     []byte("5d734c5c-1937-449a-a04d-8917b82267b5"),
		Incremental:  []byte("40cf8b78-ef1f-43e9-9253-57fc85c16ad6"),
		Fingerprint:  []byte("84fb2b2d-f6ed-4e20-bfb0-b6d79cee35b5"),
		BootId:       []byte("c81e9b68-c3fe-46af-a08c-b5fe0a5dc376"),
		AndroidId:    []byte("8b5c5fe9-a7b8-4f23-bc84-5e6730254da5"),
		Baseband:     []byte("0a73ec7a-8aac-4160-a046-8befcef3951e"),
		InnerVersion: []byte("40cf8b78-ef1f-43e9-9253-57fc85c16ad6"),
	})
	t.tlv.SetValue(bytes.NewBuffer(v))
	t.tlv.Encode(b)
}

func (t *T52D) Decode(b *bytes.Buffer) error {
	if err := t.tlv.Decode(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	deviceInfo := new(pb.DeviceInfo)
	if err := proto.Unmarshal(v.Bytes(), deviceInfo); err != nil {
		return err
	}
	panic("not implement")
}
