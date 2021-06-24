package mobileqq

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/elap5e/go-mobileqq-api/rpc"
)

func (c *Client) Auth(username, password string) error {
	uin, _ := strconv.Atoi(username)
	resp, err := c.rpc.AuthGetSessionTicketWithPassword(c.ctx, rpc.NewAuthGetSessionTicketWithPasswordRequest(uint64(uin), password))
	if err != nil {
		return err
	}
	switch resp.Code {
	case 0x00:
	case 0x02:
		if resp.CaptchaSign != "" {
			addr := "localhost:34679"
			u, _ := url.Parse(string(resp.CaptchaSign))
			captcha := resp.CaptchaSign + "\nhttp://" + addr + "/api/captcha?" + u.RawQuery
			log.Printf(">_< [info] verify captcha\n%s\n", captcha)
			type PostAuthCheckWebSignatureRequest struct {
				Uin  uint64
				Code []byte
			}
			done := make(chan *PostAuthCheckWebSignatureRequest, 1)
			go func() {
				fmt.Printf(".......... ........ >_< [info] verify captcha code: ")
				reader := bufio.NewReader(os.Stdin)
				code, _ := reader.ReadString('\n')
				done <- &PostAuthCheckWebSignatureRequest{uint64(resp.Uin), []byte(code)}
			}()
			mux := http.NewServeMux()
			mux.Handle("/api/captcha", http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						w.WriteHeader(http.StatusOK)
						fmt.Fprintln(w, tmplAuthCAPTCHA)
					case http.MethodPost:
						uin, err := strconv.Atoi(r.FormValue("uin"))
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						w.WriteHeader(http.StatusOK)
						done <- &PostAuthCheckWebSignatureRequest{uint64(uin), []byte(r.FormValue("ticket"))}
						fmt.Printf("%s\n", r.FormValue("ticket"))
					}
				},
			))
			srv := &http.Server{
				Addr:    addr,
				Handler: mux,
			}
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen:%+s\n", err)
				}
			}()
			post := <-done
			ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer func() {
				cancel()
			}()
			if err := srv.Shutdown(ctxShutDown); err != nil {
				log.Fatalf("server Shutdown Failed:%+s", err)
			}
			if resp, err = c.rpc.AuthCheckWebSignature(c.ctx, rpc.NewAuthCheckWebSignatureRequest(post.Uin, post.Code)); err != nil {
				return err
			}
		} else {
			log.Printf(">_< [info] verify picture\n\033]1337;File=name=picture.jpg;inline=1;width=11;height=2:%s\a(please check out picture.jpg)\n", base64.StdEncoding.EncodeToString(resp.PictureData))
			_ = ioutil.WriteFile("picture.jpg", resp.PictureData, 0644)
			fmt.Printf(".......... ........ >_< [info] verify picture code: ")
			reader := bufio.NewReader(os.Stdin)
			code, _ := reader.ReadString('\n')
			if resp, err = c.rpc.AuthCheckPicture(c.ctx, rpc.NewAuthCheckPictureRequest(resp.Uin, []byte(code), resp.PictureSign)); err != nil {
				return err
			}
		}
	}
	return nil
}

const tmplAuthCAPTCHA = `<!DOCTYPE html>
<html>
<head lang="zh-CN">
<meta charset="UTF-8" />
<meta name="renderer" content="webkit" />
<meta
name="viewport"
content="width=device-width,initial-scale=1,minimum-scale=1,maximum-scale=1,user-scalable=no"
/>
<title>验证码</title>
</head>
<body>
<div id="cap_iframe" style="width: 230px; height: 220px"></div>
<script type="text/javascript">
!(function () {
	var queryString = location.search
	var params = new URLSearchParams(queryString);
	var elem = document.createElement("script");
	elem.type = "text/javascript";
	elem.src = "http://captcha.qq.com/template/TCapIframeApi.js" + queryString;
	document.getElementsByTagName("head").item(0).appendChild(elem);
	elem.onload = function () {
	capInit(document.getElementById("cap_iframe"), {
		callback: function (data) {
			var xhr = new XMLHttpRequest();
			xhr.open("POST", "/api/captcha", true);
			var formData = new FormData();
			formData.append("uin", parseInt(params.get("uin")));
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
