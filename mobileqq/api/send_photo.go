package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
	"github.com/gin-gonic/gin"
)

type SendPhotoRequest struct {
	ChatID string `form:"chat_id" json:"chat_id"`
	Photo  string `form:"photo" json:"photo"`
}

func (s *Server) sendPhoto(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		botID, ok := c.Get("botID")
		if !ok {
			c.String(http.StatusUnauthorized, "error: invalid botId")
			return
		}
		req := SendPhotoRequest{}
		c.Bind(&req)
		s.handleSendPhotoRequest(ctx, botID.(string), &req, c)
	}
}

func (s *Server) handleSendPhotoRequest(ctx context.Context, botID string, req *SendPhotoRequest, c *gin.Context) {
	chatID := strings.TrimPrefix(req.ChatID, "@")
	ids := strings.Split(chatID, "_")
	_ = ids[1]
	peerID, _ := strconv.ParseUint(ids[0], 10, 64)
	// userID, _ := strconv.ParseUint(ids[1], 10, 64)
	fromID, _ := strconv.ParseUint(botID, 10, 64)

	subReq, err := s.client.NewMessageUploadImageRequest(
		req.Photo, &client.UploadRequest{
			PeerUin:   peerID,
			SecondUin: 0x0000000000000000,
			SelfUin:   fromID,
			UinType:   1,
		},
	)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		return
	}
	_, err = s.client.MessageUploadImage(ctx, botID, req.Photo, 0, subReq)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		return
	}

	s.handleSendMessageRequest(ctx, botID, &SendMessageRequest{
		ChatID: req.ChatID,
		Text: fmt.Sprintf(
			"![%s](goqq://res/image?md5=%s&type=0&uin=%d&size=%d&h=%d&w=%d)",
			subReq.FileName,
			base64.URLEncoding.EncodeToString(subReq.MD5),
			fromID,
			subReq.FileSize,
			subReq.Height,
			subReq.Width,
		),
	}, c)
}
