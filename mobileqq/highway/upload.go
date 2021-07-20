package highway

import (
	"crypto/md5"
	"io"
	"net"
	"os"
	"time"

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

	go hw.recv()

	if err := hw.Echo(); err != nil {
		return err
	}
	chunk := make([]byte, 0x00010000) // 64KiB
	offset := 0
	for {
		n, err := file.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if err := hw.uploadChunk(0x00000002, size, offset, chunk[:n], sum, ukey); err != nil {
			return err
		}
		offset += n
	}
	return nil
}

func (hw *Highway) uploadChunk(cmd uint32, size int64, offset int, body, hash, ukey []byte) error {
	sum := md5.Sum(body)
	return hw.Call(&pb.Highway_RequestHead{
		BaseHead: &pb.Highway_BaseHead{
			Version:      0x00000001,
			Uin:          hw.uin,
			Command:      "PicUp.DataUp",
			Seq:          hw.getNextSeq(),
			RetryTimes:   0x00000000, // nil
			AppId:        hw.appID,   // constant
			DataFlag:     0x00001000, // constant
			CommandId:    cmd,
			BuildVersion: "", // nil
			LocaleId:     0x00000804,
			EnvId:        0x00000000, // nil
		},
		SegmentHead: &pb.Highway_SegmentHead{
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
	}, body, nil)
}
