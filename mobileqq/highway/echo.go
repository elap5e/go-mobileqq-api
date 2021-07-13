package highway

import (
	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/pb"
)

func (hw *Highway) echo() error {
	head, err := proto.Marshal(&pb.HighwayRequestHead{
		BaseHead: &pb.HighwayBaseHead{
			Version:      0x00000001,
			Uin:          hw.uin,
			Command:      "PicUp.Echo",
			Seq:          hw.getNextSeq(),
			RetryTimes:   0x00000000,
			AppId:        hw.appID,   // constant
			DataFlag:     0x00001000, // constant
			CommandId:    0x00000000, // nil
			BuildVersion: "",         // nil
			LocaleId:     0x00000804,
			EnvId:        0x00000000, // nil
		},
	})
	if err != nil {
		return err
	}
	return hw.mustSend(head, nil)
}
