package api

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type SendPhotoRequest struct {
	ChatID                   string        `binding:"required" form:"chat_id" json:"chat_id"`
	Photo                    interface{}   ``
	Caption                  string        `form:"caption" json:"caption,omitempty"`
	ParseMode                ParseModeType `form:"parse_mode" json:"parse_mode,omitempty"`
	CaptionEntities          []interface{} `form:"caption_entities" json:"caption_entities,omitempty"`
	DisableNotification      bool          `form:"disable_notification" json:"disable_notification,omitempty"`
	ReplyToMessageID         int64         `form:"reply_to_message_id" json:"reply_to_message_id,omitempty"`
	AllowSendingWithoutReply bool          `form:"allow_sending_without_reply" json:"allow_sending_without_reply,omitempty"`
	ReplyMarkup              interface{}   `form:"reply_markup" json:"reply_markup,omitempty"`
}

type PhotoInterface interface {
	GetPhoto() interface{}
}

type PhotoString struct {
	Photo string `binding:"required" form:"photo" json:"photo"`
}

func (req *PhotoString) GetPhoto() interface{} { return req.Photo }

type PhotoInputFiles struct {
	Photo []*multipart.FileHeader `binding:"required" form:"photo" json:"photo"`
}

func (req *PhotoInputFiles) GetPhoto() interface{} { return req.Photo }

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	FileSize     int64  `json:"file_size"`
}

func (s *Server) sendPhoto(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		botID, ok := c.Get("botID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"ok":          false,
				"error_code":  http.StatusUnauthorized,
				"description": "Invalid botId",
			})
			return
		}
		req := SendPhotoRequest{}
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":          false,
				"error_code":  http.StatusBadRequest,
				"description": err.Error(),
			})
			return
		}
		var photo PhotoInterface
		photo = &PhotoString{}
		b := binding.Default(c.Request.Method, c.ContentType())
		if err := c.ShouldBindWith(photo, b); err != nil {
			photo = &PhotoInputFiles{}
			if err := c.ShouldBindWith(photo, b); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"ok":          false,
					"error_code":  http.StatusBadRequest,
					"description": err.Error(),
				})
				return
			}
		}
		req.Photo = photo.GetPhoto()

		h, err := s.handleSendPhotoRequest(ctx, botID.(string), &req, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ok":          false,
				"error_code":  http.StatusInternalServerError,
				"description": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"ok":     true,
				"result": h,
			})
		}
	}
}

func (s *Server) handleSendPhotoRequest(
	ctx context.Context,
	botID string,
	req *SendPhotoRequest,
	c *gin.Context,
) (gin.H, error) {
	peerID, _ := s.parseChatID(req.ChatID)
	fromID, _ := strconv.ParseInt(botID, 10, 64)

	fileID := ""
	subReqs, tempBlobs := []*pb.TryUploadImageRequest{}, []*client.UploadTempBlob{}
	switch photos := req.Photo.(type) {
	default:
		return nil, fmt.Errorf("Not Support")
	case string:
		for _, photo := range strings.Split(photos, "\n") {
			fileID = photo

			subReq, tempBlob, err := client.NewMessageUploadImageRequest(
				peerID, fromID, fileID, s.client.GetCacheDownloadsDir(),
			)
			if err != nil {
				return nil, err
			}

			subReqs, tempBlobs = append(subReqs, subReq), append(tempBlobs, tempBlob)
		}
	case []*multipart.FileHeader:
		for _, item := range photos {
			hash := sha256.New()
			hash.Write([]byte(item.Filename))
			uriSum := hex.EncodeToString(hash.Sum(nil))

			fileID = uriSum + "--" + path.Base(item.Filename)

			blobName := path.Join(s.client.GetCacheDownloadsDir(), fileID)
			blob, err := os.OpenFile(blobName, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				return nil, err
			}
			defer blob.Close()
			defer func() {
				if _, err := os.Stat(blobName); !os.IsNotExist(err) {
					os.Remove(blobName)
				}
			}()

			file, err := item.Open()
			if err != nil {
				return nil, err
			}
			defer file.Close()

			_, err = io.Copy(blob, file)
			if err != nil {
				return nil, err
			}

			defer func() {
				hash := sha256.New()
				hash.Write([]byte("/" + url.PathEscape(fileID)))
				uriSum := hex.EncodeToString(hash.Sum(nil))
				tempPath := path.Join(s.client.GetCacheDownloadsDir(), uriSum+".json")
				if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
					os.Remove(tempPath)
				}
			}()

			subReq, tempBlob, err := client.NewMessageUploadImageRequest(
				peerID, fromID, fileID, s.client.GetCacheDownloadsDir(),
			)
			if err != nil {
				return nil, err
			}

			subReqs, tempBlobs = append(subReqs, subReq), append(tempBlobs, tempBlob)
		}
	}

	resp, err := s.client.MessageUploadImage(ctx, botID, subReqs...)
	if err != nil {
		return nil, err
	}

	photoSizes := []PhotoSize{}

	for i, item := range resp {
		photoSizes = append(photoSizes, PhotoSize{
			FileID:       strconv.FormatInt(item.FileId, 10),
			FileUniqueID: strings.ToUpper(hex.EncodeToString(subReqs[i].GetFileMd5())),
			Width:        int64(subReqs[i].GetPictureWidth()),
			Height:       int64(subReqs[i].GetPictureHeight()),
			FileSize:     int64(subReqs[i].GetFileSize()),
		})
		for _, photoSize := range photoSizes {
			if photoSize.FileUniqueID == hex.EncodeToString(item.GetImgInfo().GetFileMd5()) {
				continue
			}
		}
		if !item.FileExist {
			addr := net.TCPAddr{
				IP:   net.IP{0, 0, 0, 0},
				Port: int(item.UploadPort[0]),
			}
			binary.LittleEndian.PutUint32(addr.IP, item.UploadIp[0])
			blobName := path.Join(s.client.GetCacheDownloadsDir(), tempBlobs[i].Name)
			hw := s.client.GetHighway(addr.String(), botID)
			if err := hw.Upload(blobName, item.UploadKey); err != nil {
				return nil, err
			}
		}
	}

	text := ""
	for _, subReq := range subReqs {
		text += fmt.Sprintf(
			"![%s](goqq://res/image?md5=%s&type=0&uin=%d&size=%d&h=%d&w=%d)",
			subReq.Filename,
			base64.URLEncoding.EncodeToString(subReq.GetFileMd5()),
			fromID,
			subReq.GetFileSize(),
			subReq.GetPictureHeight(),
			subReq.GetPictureWidth(),
		)
	}
	if req.Caption != "" {
		text += "\n" + req.Caption
	}

	h, err := s.handleSendMessageRequest(ctx, botID, &SendMessageRequest{
		ChatID: req.ChatID,
		Text:   text,
	}, c)
	if err != nil {
		return nil, err
	}
	h["photo"] = photoSizes

	return h, nil
}
