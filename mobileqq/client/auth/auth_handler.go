package auth

import (
	"context"
	"crypto/md5"
	"time"

	"github.com/elap5e/go-mobileqq-api/crypto/ecdh"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/config"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/mobileqq/rpc"
	"github.com/elap5e/go-mobileqq-api/tlv"
	"github.com/elap5e/go-mobileqq-api/util"
)

type Client interface {
	Call(
		serviceMethod string,
		c2s *codec.ClientToServerMessage,
		s2c *codec.ServerToClientMessage,
	) error
	GetNextSeq() uint32

	GetUserSignature(username string) *rpc.UserSignature
	SetUserSignature(
		ctx context.Context,
		username string,
		tlvs map[uint16]tlv.TLVCodec,
	)
	SetUserAuthSession(username string, session []byte)
	SetUserKSIDSession(username string, ksid []byte)

	SaveUserSignatures(file string) error
}

type HandlerOptions struct {
	BaseDir string
	Client  *config.ClientConfig
	Device  *config.DeviceConfig
}

type Handler struct {
	opt    *HandlerOptions
	client Client

	// crypto
	randomKey      [16]byte
	randomPassword [16]byte

	privateKey             *ecdh.PrivateKey
	serverPublicKey        *ecdh.PublicKey
	serverPublicKeyVersion uint16

	// session
	t16a, t172, t17b, t174, t402, t403 []byte

	hashedGUID     [16]byte // t401
	loginExtraData []byte   // from t537

	extraData map[uint16][]byte
}

func (h *Handler) init() {
	h.initRandomKey()
	h.initRandomPassword()
	h.initPrivateKey()
	h.initServerPublicKey()
}

func (h *Handler) GetNextSeq() uint32 { return h.client.GetNextSeq() }
func (h *Handler) GetUserSignature(username string) *rpc.UserSignature {
	return h.client.GetUserSignature(username)
}

func NewHandler(opt *HandlerOptions, client Client) *Handler {
	h := &Handler{
		opt:    opt,
		client: client,
	}
	h.init()
	return h
}

func (h *Handler) CheckCaptcha(ctx context.Context, username string, code []byte) (*Response, error) {
	return h.checkCaptchaAndGetSessionTickets(ctx,
		newCheckCaptchaAndGetSessionTicketsRequest(username, code),
	)
}

func (h *Handler) CheckPassword(ctx context.Context, username, password string) (*Response, error) {
	return h.getSessionTicketsWithPassword(ctx,
		newGetSessionTicketsWithPasswordRequest(username, password),
	)
}

func (h *Handler) CheckPicture(ctx context.Context, username string, code, sign []byte) (*Response, error) {
	return h.checkPictureAndGetSessionTickets(ctx,
		newCheckPictureAndGetSessionTicketsRequest(username, code, sign),
	)
}

func (h *Handler) CheckSMSCode(ctx context.Context, username string, code []byte) (*Response, error) {
	return h.checkSMSCodeAndGetSessionTickets(ctx,
		newCheckSMSCodeAndGetSessionTicketsRequest(username, code),
	)
}

func (h *Handler) SendSMSCode(ctx context.Context, username string) (*Response, error) {
	return h.sendSMSCode(ctx, newSendSMSCodeRequest(username))
}

func (h *Handler) SignIn(ctx context.Context, username, password string) (*Response, error) {
	sig := h.client.GetUserSignature(username)
	if len(password) != 0 {
		sig.PasswordMD5 = util.STBytesTobytes(md5.Sum([]byte(password)))
	}
	d2, ok := sig.Tickets["D2"]
	if (ok && time.Now().After(time.Unix(d2.Iss+d2.Exp, 0))) || !ok {
		return h.CheckPassword(ctx, username, password)
	} else {
		return h.getSessionTicketsWithoutPassword(ctx,
			newGetSessionTicketsWithoutPasswordRequest(username),
		)
	}
}

var handlerCtxKey struct{}

func forHandler(ctx context.Context) *Handler {
	return ctx.Value(handlerCtxKey).(*Handler)
}

func (h *Handler) withHandler(ctx context.Context) context.Context {
	return context.WithValue(ctx, handlerCtxKey, h)
}
