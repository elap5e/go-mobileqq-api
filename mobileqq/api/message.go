package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/gin-gonic/gin"
)

func (s *Server) sendMessage(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		botID, ok := c.Get("botID")
		if !ok {
			c.String(http.StatusUnauthorized, "error: invalid botId")
			return
		}

		chatID := strings.TrimPrefix(c.PostForm("chat_id"), "@")
		ids := strings.Split(chatID, "_")
		_ = ids[1]
		peerID, _ := strconv.ParseUint(ids[0], 10, 64)
		userID, _ := strconv.ParseUint(ids[1], 10, 64)
		fromID, _ := strconv.ParseUint(botID.(string), 10, 64)
		peerName := strconv.FormatUint(peerID, 10)
		userName := strconv.FormatUint(userID, 10)
		fromName := strconv.FormatUint(fromID, 10)
		text := c.PostForm("text")

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
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			return
		}
		req := client.NewMessageSendMessageRequest(
			routingHead,
			msg.GetContentHead(),
			msg.GetMessageBody(),
			0,
			nil,
		)
		resp, err := s.client.MessageSendMessage(
			ctx, botID.(string), req,
		)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			return
		}

		log.PrintMessage(
			time.Unix(resp.GetSendTime(), 0),
			peerName, userName, fromName, peerID, userID, fromID, req.GetMessageSeq(), text,
		)

		c.JSON(http.StatusOK, gin.H{
			"message_id":  req.GetMessageSeq(),
			"sender_chat": c.PostForm("chat_id"),
			"date":        resp.GetSendTime(),
		})
	}
}
