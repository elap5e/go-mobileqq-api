package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/elap5e/go-mobileqq-api/mobileqq/highway"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type Client interface {
	GetCacheDownloadsDir() string
	GetHighway(addr, username string) *highway.Highway
	GetRoutingHead(peerID, userID int64) *pb.MessageService_RoutingHead

	MessageSendMessage(ctx context.Context, username string, req *pb.MessageService_SendRequest) (*pb.MessageService_SendResponse, error)
	MessageUploadImage(ctx context.Context, username string, reqs ...*pb.Cmd0388_TryUploadImageRequest) ([]*pb.Cmd0388_TryUploadImageResponse, error)
}

type Server struct {
	ctx    context.Context
	client Client
	tokens map[string]string
}

func (s *Server) checkClient(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.client == nil {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"ok":          false,
				"error_code":  http.StatusGatewayTimeout,
				"description": "Client Not Ready",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (s *Server) checkToken(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// https://exapmle.com/bot<TOKEN>/<METHOD_NAME>
		// https://exapmle.com/bot10000:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/sendMessage
		token := c.Param("token")
		strs := strings.Split(strings.TrimPrefix(token, "bot"), ":")
		if t, ok := s.tokens[strs[0]]; !ok || token != "bot"+t {
			c.JSON(http.StatusUnauthorized, gin.H{
				"ok":          false,
				"error_code":  http.StatusUnauthorized,
				"description": "Invalid Token",
			})
			c.Abort()
			return
		}
		c.Set("botID", strs[0])
		c.Next()
	}
}

func (s *Server) parseChatID(chatID string) (peerID, userID int64) {
	chatID = strings.TrimPrefix(chatID, "@")
	ids := strings.Split(chatID, "u")
	peerID, _ = strconv.ParseInt(ids[0], 10, 64)
	if len(ids) == 2 {
		userID, _ = strconv.ParseInt(ids[1], 10, 64)
	} else {
		userID = 0
	}
	return peerID, userID
}

func NewServer(tokens map[string]string) *Server {
	return &Server{tokens: tokens}
}

func (s *Server) ResetClient(ctx context.Context, client Client) {
	// TODO: mux
	s.ctx = ctx
	s.client = client
}

func (s *Server) Run(ctx context.Context) error {
	engine := gin.Default()
	pprof.Register(engine)

	engine.Use(s.checkClient(ctx))
	engine.Use(s.checkToken(ctx))

	engine.POST("/:token/sendMessage", s.sendMessage(s.ctx))
	engine.GET("/:token/sendMessage", s.sendMessage(s.ctx))
	engine.POST("/:token/sendPhoto", s.sendPhoto(s.ctx))
	engine.GET("/:token/sendPhoto", s.sendPhoto(s.ctx))

	return engine.Run("localhost:8080")
}
