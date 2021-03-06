package tlv

import (
	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type T52D struct {
	tlv          *TLV
	deviceReport *pb.DeviceReport
}

func NewT52D(deviceReport *pb.DeviceReport) *T52D {
	return &T52D{
		tlv:          NewTLV(0x052d, 0x0000, nil),
		deviceReport: deviceReport,
	}
}

func (t *T52D) Encode(b *bytes.Buffer) {
	v, _ := proto.Marshal(t.deviceReport)
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
	deviceInfo := pb.DeviceReport{}
	if err := proto.Unmarshal(v.Bytes(), &deviceInfo); err != nil {
		return err
	}
	panic("not implement")
}
