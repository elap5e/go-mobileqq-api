package mobileqq

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
	"github.com/elap5e/go-mobileqq-api/util"
)

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
	resp *client.AuthGetSessionTicketsResponse,
) (*client.AuthGetSessionTicketsResponse, error) {
	switch resp.Code {
	case 0x00:
		return resp, nil
	case 0x02:
		if resp.CaptchaSign != "" {
			l, err := net.Listen("tcp", c.opt.Client.AuthAddress)
			if err != nil {
				log.Fatal().Err(err).
					Msg("x-x [auth] failed to start server")
			}
			addr := l.Addr().(*net.TCPAddr).String()
			u, _ := url.Parse(string(resp.CaptchaSign))
			log.Info().Msg(
				">_< [auth] verify captcha:\n" +
					"1st, legacy (deprecated): " + resp.CaptchaSign + "\n" +
					"2nd, local: http://" + addr + "/api/captcha?" + u.RawQuery,
			)
			done := make(chan string, 1)
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
						log.Info().Msg(
							">_< [auth] got captcha verify code: " + ticket,
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
					log.Fatal().Err(err).
						Msg("x-x [auth] failed to start server")
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
				log.Fatal().Err(err).
					Msg("x-x [auth] failed to shutdown server")
			}
			return c.client.AuthCheckCaptchaAndGetSessionTickets(
				c.ctx,
				client.NewAuthCheckCaptchaAndGetSessionTicketsRequest(
					resp.Username,
					[]byte(ticket),
				),
			)
		} else {
			file := fmt.Sprintf(
				"picture-%s.jpg", time.Now().Local().Format("20060102150405"),
			)
			log.Info().Msg(
				">_< [auth] picture verify, enter picture verify code:",
			)
			_ = ioutil.WriteFile(
				path.Join(c.opt.CacheDir, resp.Username, file),
				resp.PictureData,
				0600,
			)
			fmt.Println(
				"\033]1337;File=name=" + file + ";inline=1;width=11;height=2:" +
					base64.StdEncoding.EncodeToString(resp.PictureData) +
					"\a(please check out picture.jpg)",
			)
			fmt.Print(">>> ")
			code, _ := util.ReadLine(reader)
			return c.client.AuthCheckPictureAndGetSessionTickets(
				c.ctx, client.NewAuthCheckPictureAndGetSessionTicketsRequest(
					resp.Username,
					[]byte(code),
					resp.PictureSign,
				),
			)
		}
	case 0x01:
		return nil, fmt.Errorf("invalid password(0x01)")
	case 0x0a:
		return nil, fmt.Errorf("service temporarily unavailable(0x0a)")
	case 0x9a:
		return nil, fmt.Errorf("service temporarily unavailable(0x9a)")
	case 0xa0:
		log.Info().Msg(
			">_< [auth] sms verify mobile " + resp.SMSMobile +
				", enter sms verify code:",
		)
		fmt.Print(">>> ")
		code, _ := util.ReadLine(reader)
		return c.client.AuthCheckSMSAndGetSessionTickets(
			c.ctx, client.NewAuthCheckSMSAndGetSessionTicketsRequest(
				resp.Username,
				[]byte(code),
			),
		)
	case 0xa1:
		return nil, fmt.Errorf("too many sms verify requests(0xa1)")
	case 0xa2:
		return nil, fmt.Errorf("frequent sms verify requests(0xa2)")
	case 0xa4:
		return nil, fmt.Errorf("bad requests(0xa4)")
	case 0xed:
		return nil, fmt.Errorf("too many failures(0xed)")
	case 0xef:
		if resp.SMSMobile != "" {
			log.Info().Msg(
				">_< [auth] sms verify mobile " + resp.SMSMobile +
					", press ENTER to send sms verify request:",
			)
			fmt.Print(">>> ")
			_, _ = util.ReadLine(reader)
			return c.client.AuthRefreshSMSData(
				c.ctx, client.NewAuthRefreshSMSDataRequest(resp.Username),
			)
		}
	}
	return nil, fmt.Errorf("not implement(0x%02x)", resp.Code)
}

func (c *Client) Auth(username, password string) error {
	var err error
	var resp *client.AuthGetSessionTicketsResponse
	sig := c.client.GetUserSignature(username)
	if len(password) != 0 {
		sig.PasswordMD5 = util.STBytesTobytes(md5.Sum([]byte(password)))
	}
	d2, ok := sig.Tickets["D2"]
	if (ok && time.Now().After(time.Unix(d2.Iss+d2.Exp, 0))) || !ok {
		if resp, err = c.client.AuthGetSessionTicketsWithPassword(
			c.ctx,
			client.NewAuthGetSessionTicketsWithPasswordRequest(
				username,
				password,
			),
		); err != nil {
			return err
		}
	} else {
		if resp, err = c.client.AuthGetSessionTicketsWithoutPassword(
			c.ctx,
			client.NewAuthGetSessionTicketsWithoutPasswordRequest(username),
		); err != nil {
			return err
		}
	}
	for {
		if resp.Code == 0x00 {
			return nil
		}
		if resp, err = c.handleAuthResponse(resp); err != nil {
			return err
		}
	}
}
