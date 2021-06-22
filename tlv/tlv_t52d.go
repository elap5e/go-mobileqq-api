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
		Bootloader:   []byte(t.ctx.Value("Build.BOOTLOADER").(string)),
		Codename:     []byte(t.ctx.Value("Build.VERSION.CODENAME").(string)),
		Incremental:  []byte(t.ctx.Value("Build.VERSION.INCREMENTAL").(string)),
		Fingerprint:  []byte(t.ctx.Value("Build.FINGERPRINT").(string)),
		BootId:       []byte(t.ctx.Value("/proc/sys/kernel/random/boot_id").(string)),
		AndroidId:    []byte(t.ctx.Value("SystemPreferences.android_id").(string)),
		Baseband:     []byte(t.ctx.Value("android.os.SystemProperties.gsm.version.baseband").(string)),
		InnerVersion: []byte(t.ctx.Value("Build.DISPLAY|Build.VERSION.INCREMENTAL").(string)),
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
