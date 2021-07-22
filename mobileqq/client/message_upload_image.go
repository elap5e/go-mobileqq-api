package client

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/util"
)

type UploadRequest struct {
	PeerUin   int64
	SecondUin int64
	SelfUin   int64
	UinType   int
}

type UploadImageRequest struct {
	*UploadRequest
	BusinessType   int32
	Filename       string
	FileSize       int64
	Height         int32
	Width          int32
	IsContact      bool
	IsRaw          bool
	IsSnapChat     bool
	MD5            []byte
	PictureType    int32
	TransferURL    string
	TypeHotPicture int32
}

type UploadTempBlob struct {
	URL     string               `json:"url"`
	Name    string               `json:"name"`
	Size    int64                `json:"size"`
	Digests map[string][]byte    `json:"digests"`
	Photo   *UploadTempBlobPhoto `json:"photo"`
}

type UploadTempBlobPhoto struct {
	Height int32 `json:"height"`
	Width  int32 `json:"width"`
}

func newMessageUploadImageRequest(
	uri *url.URL,
	req *UploadRequest,
	cacheDir string,
) (*UploadImageRequest, *UploadTempBlob, error) {
	uri.Scheme = strings.ToLower(uri.Scheme)
	switch uri.Scheme {
	case "", "file":
		uri.Path = path.Join("/", uri.Path)
	}

	hash := sha256.New()
	hash.Write([]byte(uri.String()))
	uriSum := hex.EncodeToString(hash.Sum(nil))

	tempBlob := &UploadTempBlob{
		URL:     uri.String(),
		Name:    path.Base(uri.Path),
		Digests: make(map[string][]byte),
	}
	tempPath := path.Join(cacheDir, uriSum+".json")
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		temp, err := os.OpenFile(tempPath, os.O_RDONLY, 0644)
		if err != nil {
			return nil, nil, err
		}
		defer temp.Close()
		if err := json.NewDecoder(temp).Decode(tempBlob); err != nil {
			return nil, nil, err
		}

	} else {
		blob := &os.File{}
		switch uri.Scheme {
		case "", "file": // TODO: not safe
			blobName := path.Join(cacheDir, tempBlob.Name)
			blob, err = os.OpenFile(blobName, os.O_RDONLY, 0644)
			if err != nil {
				return nil, nil, err
			}
			defer blob.Close()

		case "http", "https":
			client := http.Client{
				CheckRedirect: func(r *http.Request, via []*http.Request) error {
					r.URL.Opaque = r.URL.Path
					return nil
				},
			}
			if httpProxy := os.Getenv("HTTP_PROXY"); httpProxy != "" {
				proxyUrl, _ := url.Parse(httpProxy)
				client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
			}
			resp, err := client.Get(tempBlob.URL)
			if err != nil {
				return nil, nil, err
			}
			defer resp.Body.Close()

			_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
			if err == nil {
				tempBlob.Name = uriSum + "--" + params["filename"]
			} else {
				tempBlob.Name = uriSum + "--" + tempBlob.Name
			}

			tempExt := path.Ext(tempBlob.Name)
			if tempExt == "" {
				switch resp.Header.Get("Content-Type") {
				case "image/jpeg", "image/jpg":
					tempBlob.Name += ".jpg"
				case "image/png":
					tempBlob.Name += ".png"
				case "image/bmp":
					tempBlob.Name += ".bmp"
				case "image/gif":
					tempBlob.Name += ".gif"
				case "image/heic":
					tempBlob.Name += ".heic"
				case "image/heif":
					tempBlob.Name += ".heif"
				default:
					tempBlob.Name += ".gif" // TODO: fix
				}
			}

			blobName := path.Join(cacheDir, tempBlob.Name)
			blob, err = os.OpenFile(blobName, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				return nil, nil, err
			}
			defer blob.Close()

			_, err = io.Copy(blob, resp.Body)
			if err != nil {
				return nil, nil, err
			}
		}

		blob.Seek(0, io.SeekStart)
		im, _, err := image.DecodeConfig(blob)
		if err != nil {
			return nil, nil, err
		}

		hash = md5.New()
		blob.Seek(0, io.SeekStart)
		n, err := io.Copy(hash, blob)
		if err != nil {
			return nil, nil, err
		}
		tempBlob.Size = n
		tempBlob.Digests["md5"] = hash.Sum(nil)

		hash = sha256.New()
		blob.Seek(0, io.SeekStart)
		n, err = io.Copy(hash, blob)
		if err != nil {
			return nil, nil, err
		}
		tempBlob.Size = n
		tempBlob.Digests["sha256"] = hash.Sum(nil)

		tempBlob.Photo = &UploadTempBlobPhoto{
			Height: int32(im.Height),
			Width:  int32(im.Width),
		}

		temp, err := os.OpenFile(tempPath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, nil, err
		}
		defer temp.Close()

		enc := json.NewEncoder(temp)
		enc.SetIndent("", "  ")
		if err := enc.Encode(&tempBlob); err != nil {
			return nil, nil, err
		}
	}

	ext := path.Ext(tempBlob.Name)
	photo := tempBlob.Photo
	if photo == nil {
		photo = &UploadTempBlobPhoto{}
	}

	return &UploadImageRequest{
		UploadRequest:  req,
		BusinessType:   0x000003ee, // 1006
		Filename:       util.HashToString(tempBlob.Digests["md5"]) + ext,
		FileSize:       tempBlob.Size,
		Height:         photo.Height,
		Width:          photo.Width,
		IsContact:      false, // nil
		IsRaw:          true,
		IsSnapChat:     false, // nil
		MD5:            tempBlob.Digests["md5"],
		PictureType:    util.ParseExtToPictureType(ext),
		TransferURL:    "",         // nil
		TypeHotPicture: 0x00000000, // nil
	}, tempBlob, nil
}

func NewMessageUploadImageRequest(
	peerUin int64,
	selfUin int64,
	fileID string,
	cacheDir string,
) (*pb.Cmd0388_TryUploadImageRequest, *UploadTempBlob, error) {
	uri, err := url.Parse(fileID)
	if err != nil {
		return nil, nil, err
	}

	req, tempBlob, err := newMessageUploadImageRequest(
		uri, &UploadRequest{
			PeerUin:   peerUin,
			SecondUin: 0x0000000000000000,
			SelfUin:   selfUin,
			UinType:   1,
		}, cacheDir,
	)
	if err != nil {
		return nil, nil, err
	}

	buType := int32(1)
	if req.UinType != 1 {
		buType = 2
	}
	originalPicture := int32(0)
	if req.IsRaw {
		originalPicture = 1
	}

	return &pb.Cmd0388_TryUploadImageRequest{
		GroupCode:       req.PeerUin,
		SrcUin:          req.SelfUin,
		FileId:          0x0000000000000000, // nil
		FileMd5:         req.MD5,
		FileSize:        req.FileSize,
		Filename:        req.Filename,
		SrcTerm:         0x00000005,
		PlatformType:    0x00000009,
		BuType:          buType,
		PictureWidth:    req.Width,
		PictureHeight:   req.Height,
		PictureType:     req.PictureType,
		BuildVersion:    "",         // placeholder
		InnerIp:         0x00000000, // nil
		AppPictureType:  req.BusinessType,
		OriginalPicture: originalPicture,
		FileIndex:       nil,                // nil
		DstUin:          0x0000000000000000, // nil
		ServerUpload:    req.TypeHotPicture,
		TransferUrl:     req.TransferURL,
	}, tempBlob, nil
}

func (c *Client) MessageUploadImage(
	ctx context.Context,
	username string,
	reqs ...*pb.Cmd0388_TryUploadImageRequest,
) ([]*pb.Cmd0388_TryUploadImageResponse, error) {
	for _, req := range reqs {
		req.BuildVersion = c.cfg.Client.Revision
	}

	buf, err := proto.Marshal(&pb.Cmd0388_Request{
		NetType:        0x00000003,
		SubCmd:         0x00000001,
		TryUploadImage: reqs,
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
	resp := pb.Cmd0388_Response{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}

	util.DumpServerToClientMessage(&s2c, &resp)
	return resp.GetTryUploadImage(), nil
}
