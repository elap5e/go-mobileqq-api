package highway

import (
	"github.com/elap5e/go-mobileqq-api/pb"
)

func (hw *Highway) Echo() error {
	return hw.Call(&pb.Highway_RequestHead{
		BaseHead: &pb.Highway_BaseHead{
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
	}, nil, nil)
}
