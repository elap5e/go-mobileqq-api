package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type Client interface {
	MessageSendMessage(ctx context.Context, username string, req *pb.MessageSendMessageRequest) (*pb.MessageSendMessageResponse, error)
	MessageUploadImage(ctx context.Context, username string, name string, fileID uint64, req *client.UploadImageRequest) (*pb.Cmd0X0388Response, error)

	NewMessageUploadImageRequest(name string, req *client.UploadRequest) (*client.UploadImageRequest, error)
}

type Server struct {
	client Client
	tokens map[string]string
}

func (s *Server) checkToken(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// https://exapmle.com/bot<TOKEN>/<METHOD_NAME>
		// https://exapmle.com/bot10000:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/sendMessage
		token := c.Param("token")
		strs := strings.Split(strings.TrimPrefix(token, "bot"), ":")
		if t, ok := s.tokens[strs[0]]; !ok || token != "bot"+t {
			c.String(http.StatusUnauthorized, "error: invalid token")
			c.Abort()
			return
		}
		c.Set("botID", strs[0])
		c.Next()
	}
}

func NewServer(client Client, tokens map[string]string) *Server {
	return &Server{client: client, tokens: tokens}
}

func (s *Server) Run(ctx context.Context) error {
	engine := gin.Default()
	pprof.Register(engine)

	engine.Use(s.checkToken(ctx))

	engine.POST("/:token/sendMessage", s.sendMessage(ctx))
	engine.GET("/:token/sendMessage", s.sendMessage(ctx))
	engine.POST("/:token/sendPhoto", s.sendPhoto(ctx))
	engine.GET("/:token/sendPhoto", s.sendPhoto(ctx))

	return engine.Run(":8080")
}
