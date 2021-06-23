package mobileqq

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

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
	var elem = document.createElement("script");
	elem.type = "text/javascript";
	elem.src = "http://captcha.qq.com/template/TCapIframeApi.js" + location.search;
	document.getElementsByTagName("head").item(0).appendChild(elem);
	elem.onload = function () {
	capInit(document.getElementById("cap_iframe"), {
		callback: function (data) {
			var xhr = new XMLHttpRequest();
			xhr.open("POST", "/api/captcha", true);
			xhr.setRequestHeader(
				"Content-Type",
				"application/json; charset=UTF-8"
			);
			xhr.onload = function (e) { window.close(); };
			xhr.send(JSON.stringify(data));
		},
		showHeader: !1,
	});
	};
})();
</script>
</body>
</html>`

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		want *Client
	}{
		{
			name: "0x00",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()
			c.Auth("123456", "123456")
			http.HandleFunc("/auth/captcha.html", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, tmplAuthCAPTCHA)
			}))
			http.HandleFunc("/api/captcha", func(response http.ResponseWriter, request *http.Request) {
				body, _ := ioutil.ReadAll(request.Body)
				log.Println(string(body))
			})
			http.ListenAndServe(":8080", nil)
		})
	}
}
