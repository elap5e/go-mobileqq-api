package highway

import (
	"crypto/md5"
	"io"
	"net"
	"os"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/pb"
)

func (hw *Highway) Upload(name string, ukey []byte) error {
	file, err := os.OpenFile(name, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := md5.New()
	size, err := io.Copy(hash, file)
	if err != nil {
		return err
	}
	file.Seek(0, io.SeekStart)
	sum := hash.Sum(nil)

	hw.conn, err = net.DialTimeout("tcp", hw.addr, 30*time.Second)
	if err != nil {
		return err
	}
	defer hw.conn.Close()

	if err := hw.echo(); err != nil {
		return err
	}
	offset := 0
	chunk := make([]byte, 0x00010000) // 64KB
	for {
		n, err := file.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if err := hw.uploadChunk(size, offset, chunk[:n], sum, ukey); err != nil {
			return err
		}
		offset += n
	}
	return nil
}

func (hw *Highway) uploadChunk(size int64, offset int, body, hash, ukey []byte) error {
	sum := md5.Sum(body)
	req, err := proto.Marshal(&pb.HighwayRequestHead{
		BaseHead: &pb.HighwayBaseHead{
			Version:      0x00000001,
			Uin:          hw.uin,
			Command:      "PicUp.DataUp",
			Seq:          hw.getNextSeq(),
			RetryTimes:   0x00000000, // nil
			AppId:        hw.appID,   // constant
			DataFlag:     0x00001000, // constant
			CommandId:    0x00000002, // TODO: fix
			BuildVersion: "",         // nil
			LocaleId:     0x00000804,
			EnvId:        0x00000000, // nil
		},
		SegmentHead: &pb.HighwaySegmentHead{
			ServiceId:     0x00000000, // nil
			FileSize:      uint64(size),
			DataOffset:    uint64(offset),
			DataLength:    uint32(len(body)),
			ReturnCode:    0x00000000, // nil
			ServiceTicket: ukey,
			Flag:          0x00000000, // nil
			Md5:           sum[:],
			FileMd5:       hash,
			CacheAddress:  0x00000000, // nil
			QueryTimes:    0x00000000, // nil
			UpdateCacheIp: 0x00000000, // nil
			CachePort:     0x00000000, // nil
		},
		ExtendInfo: []byte{},
	})
	if err != nil {
		return err
	}
	return hw.mustSend(req, body)
}
