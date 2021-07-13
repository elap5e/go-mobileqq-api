package client

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/mobileqq/highway"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type UploadRequest struct {
	PeerUin   uint64
	SecondUin uint64
	SelfUin   uint64
	UinType   int
}

type UploadImageRequest struct {
	*UploadRequest
	BusinessType   uint32
	FileName       string
	FileSize       uint64
	Height         uint32
	Width          uint32
	IsContact      bool
	IsRaw          bool
	IsSnapChat     bool
	MD5            []byte
	PictureType    uint32
	TransferURL    string
	TypeHotPicture uint32
}

type UploadTempBlob struct {
	Path   string `json:"path"`
	Name   string `json:"name"`
	Hash   []byte `json:"hash"`
	Size   uint64 `json:"size"`
	Height uint32 `json:"height"`
	Width  uint32 `json:"width"`
}

func ParseExt2PictureType(ext string) uint32 {
	switch ext {
	case ".jpg", ".jepg":
		return 1000
	case ".png":
		return 1001
	case ".webp":
		return 1002
	case ".sharpp":
		return 1004
	case ".bmp":
		return 1005
	case ".gif":
		return 2000
	case ".apng":
		return 2001
	}
	return 0
}

func (c *Client) NewMessageUploadImageRequest(
	name string,
	req *UploadRequest,
) (*UploadImageRequest, error) {
	rich, err := url.Parse(name)
	if err != nil {
		return nil, err
	}

	tempBlob := UploadTempBlob{Path: name}
	tempName := ""

	switch strings.ToLower(rich.Scheme) {
	case "http", "https":
		hash := sha256.New()
		hash.Write([]byte(name))
		sum := hex.EncodeToString(hash.Sum(nil))

		tempBlob = UploadTempBlob{Path: name}
		tempName = path.Join(c.cfg.CacheDir, "downloads", sum+".json")
		_, err := os.Stat(tempName)
		if !os.IsNotExist(err) {
			temp, err := os.OpenFile(tempName, os.O_RDONLY, 0644)
			if err != nil {
				return nil, err
			}
			defer temp.Close()
			if err := json.NewDecoder(temp).Decode(&tempBlob); err != nil {
				return nil, err
			}
			ext := path.Ext(tempBlob.Name)
			return &UploadImageRequest{
				UploadRequest:  req,
				BusinessType:   0x000003ee, // 1006
				FileName:       strings.ToUpper(hex.EncodeToString(tempBlob.Hash)) + ext,
				FileSize:       tempBlob.Size,
				Height:         tempBlob.Height,
				Width:          tempBlob.Width,
				IsContact:      false, // nil
				IsRaw:          true,
				IsSnapChat:     false, // nil
				MD5:            tempBlob.Hash,
				PictureType:    ParseExt2PictureType(ext),
				TransferURL:    "",         // nil
				TypeHotPicture: 0x00000000, // nil
			}, nil
		}

		proxyUrl, _ := url.Parse(os.Getenv("HTTP_PROXY"))
		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		}
		resp, err := client.Get(name)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if err == nil {
			name = params["filename"]
		} else {
			name = path.Base(name)
		}

		tempBlob.Name = sum + "--" + name
		name = path.Join(c.cfg.CacheDir, "downloads", tempBlob.Name)
		blob, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		defer blob.Close()

		_, err = io.Copy(blob, resp.Body)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(name, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	im, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}
	file.Seek(0, io.SeekStart)

	hash := md5.New()
	n, err := io.Copy(hash, file)
	if err != nil {
		return nil, err
	}
	file.Seek(0, io.SeekStart)
	sum := hash.Sum(nil)

	ext := path.Ext(name)

	switch strings.ToLower(rich.Scheme) {
	case "http", "https":
		tempBlob.Hash = sum
		tempBlob.Size = uint64(n)
		tempBlob.Height = uint32(im.Height)
		tempBlob.Width = uint32(im.Width)
		temp, err := os.OpenFile(tempName, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		defer temp.Close()
		if err := json.NewEncoder(temp).Encode(&tempBlob); err != nil {
			return nil, err
		}
	}

	return &UploadImageRequest{
		UploadRequest:  req,
		BusinessType:   0x000003ee, // 1006
		FileName:       strings.ToUpper(hex.EncodeToString(sum)) + ext,
		FileSize:       uint64(n),
		Height:         uint32(im.Height),
		Width:          uint32(im.Width),
		IsContact:      false, // nil
		IsRaw:          true,
		IsSnapChat:     false, // nil
		MD5:            sum,
		PictureType:    ParseExt2PictureType(ext),
		TransferURL:    "",         // nil
		TypeHotPicture: 0x00000000, // nil
	}, nil
}

func (c *Client) MessageUploadImage(
	ctx context.Context,
	username string,
	name string,
	fileID uint64,
	req *UploadImageRequest,
) (*pb.Cmd0X0388Response, error) {
	buType := uint32(1)
	if req.UinType != 1 {
		buType = 2
	}
	originalPicture := uint32(0)
	if req.IsRaw {
		originalPicture = 1
	}

	buf, err := proto.Marshal(&pb.Cmd0X0388Request{
		NetType: 0x00000003,
		SubCmd:  0x00000001,
		TryUploadImage: []*pb.TryUploadImageRequest{{
			GroupCode:       req.PeerUin,
			SrcUin:          req.SelfUin,
			FileId:          fileID,
			FileMd5:         req.MD5,
			FileSize:        req.FileSize,
			FileName:        req.FileName,
			SrcTerm:         0x00000005,
			PlatformType:    0x00000009,
			BuType:          buType,
			PictureWidth:    req.Width,
			PictureHeight:   req.Height,
			PictureType:     req.PictureType,
			BuildVersion:    c.cfg.Client.Revision,
			InnerIp:         0x00000000, // nil
			AppPictureType:  req.BusinessType,
			OriginalPicture: originalPicture,
			FileIndex:       nil,                // nil
			DstUin:          0x0000000000000000, // nil
			ServerUpload:    req.TypeHotPicture,
			TransferUrl:     req.TransferURL,
		}},
	})
	if err != nil {
		return nil, err
	}
	c2s, s2c := codec.ClientToServerMessage{
		Username: username,
		Buffer:   buf,
		Simple:   true,
	}, codec.ServerToClientMessage{}
	err = c.rpc.Call(ServiceMethodMessageUploadImageGroup, &c2s, &s2c)
	if err != nil {
		return nil, err
	}
	resp := pb.Cmd0X0388Response{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}

	dumpServerToClientMessage(&s2c, &resp)

	for _, item := range resp.TryUploadImage {
		rich, err := url.Parse(name)
		if err != nil {
			return nil, err
		}
		switch strings.ToLower(rich.Scheme) {
		case "http", "https":
			hash := sha256.New()
			hash.Write([]byte(name))
			sum := hex.EncodeToString(hash.Sum(nil))
			tempName := path.Join(c.cfg.CacheDir, "downloads", sum+".json")
			temp, err := os.OpenFile(tempName, os.O_RDONLY, 0644)
			if err != nil {
				return nil, err
			}
			defer temp.Close()
			tempBlob := UploadTempBlob{}
			if err := json.NewDecoder(temp).Decode(&tempBlob); err != nil {
				return nil, err
			}
			name = path.Join(c.cfg.CacheDir, "downloads", tempBlob.Name)
			defer func() {
				_, err := os.Stat(name)
				if !os.IsNotExist(err) {
					os.Remove(name)
				}
			}()
		}
		if !item.FileExist {
			addr := net.TCPAddr{
				IP:   net.IP{0, 0, 0, 0},
				Port: int(item.UploadPort[0]),
			}
			binary.LittleEndian.PutUint32(addr.IP, item.UploadIp[0])
			hw := highway.NewHighway(
				addr.String(),
				strconv.FormatUint(req.SelfUin, 10),
				c.cfg.Client.AppID,
			)
			if err := hw.Upload(name, item.UploadKey); err != nil {
				return nil, err
			}
		}
	}

	return &resp, nil
}
