package mobileqq

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/elap5e/go-mobileqq-api/rpc"
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

func (c *Client) handleAuthResponse(resp *rpc.AuthGetSessionTicketResponse) (*rpc.AuthGetSessionTicketResponse, error) {
	switch resp.Code {
	case 0x00:
		return resp, nil
	case 0x02:
		if resp.CaptchaSign != "" {
			l, err := net.Listen("tcp", ":0")
			if err != nil {
				log.Fatalf("listen:%+s\n", err)
			}
			addr := l.Addr().(*net.TCPAddr).String()
			u, _ := url.Parse(string(resp.CaptchaSign))
			captcha := ".......... ........ >_< [info] 1st, submit captcha manually: " + resp.CaptchaSign + "\n"
			captcha += ".......... ........ >_< [info] 2nd, submit captcha by popup: http://" + addr + "/api/captcha?" + u.RawQuery
			log.Printf(">_< [info] verify captcha:\n%s\n", captcha)
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
						log.Printf(">_< [info] got captcha verify code: %s", ticket)
						done <- ticket
					}
				},
			))
			srv := &http.Server{
				Handler: mux,
			}
			go func() {
				if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen:%+s\n", err)
				}
			}()
			ticket := <-done
			ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer func() {
				cancel()
			}()
			if err := srv.Shutdown(ctxShutDown); err != nil {
				log.Fatalf("server Shutdown Failed:%+s", err)
			}
			return c.rpc.AuthCheckCaptchaAndGetSessionTicket(c.ctx, rpc.NewAuthCheckCaptchaAndGetSessionTicketRequest(resp.Uin, []byte(ticket)))
		} else {
			log.Printf(">_< [info] picture verify:\n\033]1337;File=name=picture.jpg;inline=1;width=11;height=2:%s\a(please check out picture.jpg)\n", base64.StdEncoding.EncodeToString(resp.PictureData))
			_ = ioutil.WriteFile("picture.jpg", resp.PictureData, 0644)
			fmt.Printf(".......... ........ >_< [info] picture verify code: ")
			code, _ := util.ReadLine(reader)
			return c.rpc.AuthCheckPictureAndGetSessionTicket(c.ctx, rpc.NewAuthCheckPictureAndGetSessionTicketRequest(resp.Uin, []byte(code), resp.PictureSign))
		}
	case 0x01:
		return nil, fmt.Errorf("invalid password(0x01)")
	case 0x9a:
		return nil, fmt.Errorf("service temporarily unavailable(0x9a)")
	case 0xa0:
		fmt.Printf(".......... ........ >_< [info] sms mobile verify code: ")
		code, _ := util.ReadLine(reader)
		return c.rpc.AuthCheckSMSAndGetSessionTicket(c.ctx, rpc.NewAuthCheckSMSAndGetSessionTicketRequest(resp.Uin, []byte(code)))
	case 0xa1:
		return nil, fmt.Errorf("too many sms verify requests(0xa1)")
	case 0xa2:
		return nil, fmt.Errorf("frequent sms verify requests(0xa2)")
	case 0xed:
		return nil, fmt.Errorf("too many failures(0xed)")
	case 0xef:
		if resp.SMSMobile != "" {
			log.Printf(">_< [info] verify sms mobile %s", resp.SMSMobile)
			fmt.Printf(".......... ........ >_< [info] press ENTER to send sms mobile verify request: ")
			_, _ = util.ReadLine(reader)
			return c.rpc.AuthRefreshSMSData(c.ctx, rpc.NewAuthRefreshSMSDataRequest(resp.Uin))
		}
	}
	return nil, fmt.Errorf("not implement(0x%02x)", resp.Code)
}

func (c *Client) Auth(username, password string) error {
	uin, _ := strconv.Atoi(username)
	resp, err := c.rpc.AuthGetSessionTicketWithPassword(c.ctx, rpc.NewAuthGetSessionTicketWithPasswordRequest(uint64(uin), password))
	if err != nil {
		return err
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
