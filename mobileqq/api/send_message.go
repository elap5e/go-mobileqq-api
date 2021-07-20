package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
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

		h, err := s.handleSendMessageRequest(ctx, botID.(string), &req, c)
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

func (s *Server) handleSendMessageRequest(
	ctx context.Context,
	botID string,
	req *SendMessageRequest,
	c *gin.Context,
) (gin.H, error) {
	peerID, userID := s.parseChatID(req.ChatID)
	fromID, _ := strconv.ParseInt(botID, 10, 64)
	text := req.Text

	elems, err := mark.NewDecoder(peerID, userID, fromID).
		Decode([]byte(text))
	if err != nil {
		return nil, err
	}
	msg := pb.MessageCommon_Message{
		MessageBody: &pb.IMMessageBody_MessageBody{
			RichText: &pb.IMMessageBody_RichText{
				Elements: elems,
			},
		},
	}
	subReq := client.NewMessageSendMessageRequest(
		s.client.GetRoutingHead(peerID, userID),
		msg.GetContentHead(),
		msg.GetMessageBody(),
		0,
		nil,
	)
	resp, err := s.client.MessageSendMessage(ctx, botID, subReq)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"message_id":  subReq.GetMessageSeq(),
		"sender_chat": c.PostForm("chat_id"),
		"date":        resp.GetSendTime(),
		"text":        text,
	}, nil
}
