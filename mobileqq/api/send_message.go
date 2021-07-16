package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type ParseModeType string

const (
	ParseModeHTML       ParseModeType = "HTML"
	ParseModeMarkdown   ParseModeType = "Markdown"
	ParseModeMarkdownV2 ParseModeType = "MarkdownV2"
)

type SendMessageRequest struct {
	ChatID                   string        `binding:"required" form:"chat_id" json:"chat_id"`
	Text                     string        `binding:"required" form:"text" json:"text"`
	ParseMode                ParseModeType `form:"parse_mode" json:"parse_mode,omitempty"`
	Entities                 []interface{} `form:"entities" json:"entities,omitempty"`
	DisableWebPagePreview    bool          `form:"disable_web_page_preview" json:"disable_web_page_preview,omitempty"`
	DisableNotification      bool          `form:"disable_notification" json:"disable_notification,omitempty"`
	ReplyToMessageID         int64         `form:"reply_to_message_id" json:"reply_to_message_id,omitempty"`
	AllowSendingWithoutReply bool          `form:"allow_sending_without_reply" json:"allow_sending_without_reply,omitempty"`
	ReplyMarkup              interface{}   `form:"reply_markup" json:"reply_markup,omitempty"`
}

func (s *Server) sendMessage(ctx context.Context) gin.HandlerFunc {
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
		req := SendMessageRequest{}
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":          false,
				"error_code":  http.StatusBadRequest,
				"description": err.Error(),
			})
			return
		}
		s.handleSendMessageRequest(ctx, botID.(string), &req, c)
	}
}

func (s *Server) handleSendMessageRequest(ctx context.Context, botID string, req *SendMessageRequest, c *gin.Context) {
	peerID, userID := s.parseChatID(req.ChatID)
	fromID, _ := strconv.ParseUint(botID, 10, 64)
	peerName := strconv.FormatUint(peerID, 10)
	userName := strconv.FormatUint(userID, 10)
	fromName := strconv.FormatUint(fromID, 10)
	text := req.Text

	routingHead := &pb.RoutingHead{}
	if peerID == 0 {
		routingHead = &pb.RoutingHead{C2C: &pb.C2C{ToUin: userID}}
	} else if userID == 0 {
		routingHead = &pb.RoutingHead{Group: &pb.Group{Code: peerID}}
	} else {
		routingHead = &pb.RoutingHead{
			GroupTemp: &pb.GroupTemp{Uin: peerID, ToUin: userID},
		}
	}

	msg := pb.Message{}
	if err := mark.Unmarshal([]byte(text), &msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok":          false,
			"error_code":  http.StatusInternalServerError,
			"description": err.Error(),
		})
		return
	}
	subReq := client.NewMessageSendMessageRequest(
		routingHead,
		msg.GetContentHead(),
		msg.GetMessageBody(),
		0,
		nil,
	)
	resp, err := s.client.MessageSendMessage(ctx, botID, subReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok":          false,
			"error_code":  http.StatusInternalServerError,
			"description": err.Error(),
		})
		return
	}

	log.PrintMessage(
		time.Unix(resp.GetSendTime(), 0),
		peerName, userName, fromName, peerID, userID, fromID, subReq.GetMessageSeq(), text,
	)

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"result": gin.H{
			"message_id":  subReq.GetMessageSeq(),
			"sender_chat": c.PostForm("chat_id"),
			"date":        resp.GetSendTime(),
		},
	})
}
