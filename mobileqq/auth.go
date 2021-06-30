package mobileqq

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/elap5e/go-mobileqq-api/rpc"
	"github.com/elap5e/go-mobileqq-api/util"
)

const PATH_TO_AUTH_PICTURE_JPG = "picture.jpg"

const tmplAuthCaptcha = `<!DOCTYPE html>
<html>
<head lang="zh-CN">
	<meta charset="UTF-8" />
	<meta name="renderer" content="webkit" />
	<meta name="viewport" content="width=device-width,initial-scale=1,minimum-scale=1,maximum-scale=1,user-scalable=no" />
	<title>验证码</title>
</head>
<body>
	<div id="cap_iframe" style="width: 230px; height: 220px"></div>
	<script type="text/javascript">
		!(function () {
			var elem = document.createElement("script");
			elem.type = "text/javascript";
			elem.src = "http://captcha.qq.com/template/TCapIframeApi.js" + location.search;
			document.getElementsByTagName("head").item(0).appendChild(elem);
			elem.onload = function () {
				capInit(document.getElementById("cap_iframe"), {
					callback: function (data) {
						var xhr = new XMLHttpRequest();
						xhr.open("POST", "/api/captcha", true);
						var formData = new FormData();
						formData.append("ticket", data.ticket);
						formData.append("code", data.ret);
						xhr.onload = function (e) { window.close(); };
						xhr.send(formData);
					},
					showHeader: !1,
				});
			};
		})();
	</script>
</body>
</html>`

var reader = bufio.NewReader(os.Stdin)

func (c *Client) handleAuthResponse(
	resp *rpc.AuthGetSessionTicketsResponse,
) (*rpc.AuthGetSessionTicketsResponse, error) {
	switch resp.Code {
	case 0x00:
		return resp, nil
	case 0x02:
		if resp.CaptchaSign != "" {
			l, err := net.Listen("tcp", c.cfg.Client.AuthAddress)
			if err != nil {
				log.Fatalf("listen:%+s\n", err)
			}
			addr := l.Addr().(*net.TCPAddr).String()
			u, _ := url.Parse(string(resp.CaptchaSign))
			info := ".......... ........ >_< [info] 1st, submit captcha manually: " + resp.CaptchaSign + "\n"
			info += ".......... ........ >_< [info] 2nd, submit captcha by popup: http://" + addr + "/api/captcha?" + u.RawQuery
			log.Printf(">_< [info] verify captcha:\n%s\n", info)
			done := make(chan string, 1)
			// ctx, cancel := context.WithCancel(context.Background())
			// go func() {
			// 	fmt.Printf(".......... ........ >_< [info] captcha verify code: ")
			// 	ticket, _ := util.ReadLineWithCtx(ctx, reader)
			// 	done <- ticket
			// }()
			mux := http.NewServeMux()
			mux.Handle("/api/captcha", http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						w.WriteHeader(http.StatusOK)
						fmt.Fprintln(w, tmplAuthCaptcha)
					case http.MethodPost:
						w.WriteHeader(http.StatusOK)
						ticket := r.FormValue("ticket")
						// fmt.Println()
						// cancel()
						log.Printf(
							">_< [info] got captcha verify code: %s",
							ticket,
						)
						done <- ticket
					}
				},
			))
			srv := &http.Server{
				Handler: mux,
			}
			go func() {
				err := srv.Serve(l)
				if err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen:%+s\n", err)
				}
			}()
			ticket := <-done
			ctxShutDown, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer func() {
				cancel()
			}()
			if err := srv.Shutdown(ctxShutDown); err != nil {
				log.Fatalf("server Shutdown Failed:%+s", err)
			}
			return c.rpc.AuthCheckCaptchaAndGetSessionTickets(
				c.ctx,
				rpc.NewAuthCheckCaptchaAndGetSessionTicketsRequest(
					resp.Username,
					[]byte(ticket),
				),
			)
		} else {
			info := ".......... ........ >_< [info] picture verify code: "
			log.Printf(
				">_< [info] picture verify:\n\033]1337;File=name=picture.jpg;inline=1;width=11;height=2:%s\a(please check out picture.jpg)\n%s",
				base64.StdEncoding.EncodeToString(resp.PictureData), info,
			)
			_ = ioutil.WriteFile(
				path.Join(
					c.cfg.CacheDir,
					fmt.Sprintf(
						"%s-picture-%s.jpg",
						resp.Username,
						time.Now().Local().Format("20060102150405"),
					),
				),
				resp.PictureData,
				0600,
			)
			code, _ := util.ReadLine(reader)
			return c.rpc.AuthCheckPictureAndGetSessionTickets(
				c.ctx,
				rpc.NewAuthCheckPictureAndGetSessionTicketsRequest(
					resp.Username,
					[]byte(code),
					resp.PictureSign,
				),
			)
		}
	case 0x01:
		return nil, fmt.Errorf("invalid password(0x01)")
	case 0x0a: // TODO: check
		return nil, fmt.Errorf("service temporarily unavailable(0x0a)")
	case 0x9a:
		return nil, fmt.Errorf("service temporarily unavailable(0x9a)")
	case 0xa0: // TODO: check
		info := ".......... ........ >_< [info] sms verification code: "
		log.Printf(
			">_< [info] sms verification mobile %s\n%s",
			resp.SMSMobile, info,
		)
		code, _ := util.ReadLine(reader)
		return c.rpc.AuthCheckSMSAndGetSessionTickets(
			c.ctx,
			rpc.NewAuthCheckSMSAndGetSessionTicketsRequest(
				resp.Username,
				[]byte(code),
			),
		)
	case 0xa1:
		return nil, fmt.Errorf("too many sms verification requests(0xa1)")
	case 0xa2:
		return nil, fmt.Errorf("frequent sms verification requests(0xa2)")
	case 0xa4:
		return nil, fmt.Errorf("bad requests(0xa4)")
	case 0xed:
		return nil, fmt.Errorf("too many failures(0xed)")
	case 0xef:
		if resp.SMSMobile != "" {
			info := ".......... ........ >_< [info] press ENTER to send sms verification request: "
			log.Printf(
				">_< [info] sms verification mobile %s\n%s",
				resp.SMSMobile, info,
			)
			_, _ = util.ReadLine(reader)
			return c.rpc.AuthRefreshSMSData(
				c.ctx,
				rpc.NewAuthRefreshSMSDataRequest(resp.Username),
			)
		}
	}
	return nil, fmt.Errorf("not implement(0x%02x)", resp.Code)
}

func (c *Client) Auth(username, password string) error {
	var err error
	var resp *rpc.AuthGetSessionTicketsResponse
	sig := c.rpc.GetUserSignature(username)
	if len(password) != 0 {
		sig.PasswordMD5 = util.STBytesTobytes(md5.Sum([]byte(password)))
	}
	d2, ok := sig.Tickets["D2"]
	if (ok && time.Now().After(time.Unix(d2.Iss+d2.Exp, 0))) || !ok {
		if resp, err = c.rpc.AuthGetSessionTicketsWithPassword(
			c.ctx,
			rpc.NewAuthGetSessionTicketsWithPasswordRequest(
				username,
				password,
			),
		); err != nil {
			return err
		}
	} else {
		if resp, err = c.rpc.AuthGetSessionTicketsWithoutPassword(
			c.ctx,
			rpc.NewAuthGetSessionTicketsWithoutPasswordRequest(username),
		); err != nil {
			return err
		}
	}
	for {
		if resp.Code == 0x00 {
			tresp, err := c.rpc.AccountUpdateStatus(
				c.ctx,
				rpc.NewAccountUpdateStatusRequest(
					resp.Uin,
					rpc.PushRegisterInfoStatusOnline,
					false,
				),
			)
			if err != nil {
				return err
			}
			jresp, _ := json.MarshalIndent(tresp, "", "  ")
			log.Printf("AccountUpdateStatusResponse\n%s", jresp)
			return nil
		}
		if resp, err = c.handleAuthResponse(resp); err != nil {
			return err
		}
	}
}
