package proxy

import (
	"encoding/json"
	"ipgw/api/global"
	. "ipgw/base"
	. "ipgw/core/global"
	"ipgw/ctx"
	"net/http"
	"strings"
)

// 使用账号密码代理请求
func proxyWithUP(c *ctx.Ctx, url, method, headers, body string) {
	// 先登录
	global.LoginWithUP(c, url, isWebvpn(url))

	// 代理请求
	proxy(c, url, method, headers, body)
}

func proxyWithC(c *ctx.Ctx, url, method, headers, body string) {
	// 先登录
	global.LoginWithC(c, url, isWebvpn(url))

	// 代理请求
	proxy(c, url, method, headers, body)
}

func proxy(c *ctx.Ctx, url, method, headers, body string) {
	// 登陆成功，构造请求
	req := buildRequest(url, method, headers, body)

	// 访问
	global.SendRequest(c, req)

	Info(ReadBody(c.Response))
}

// 构造请求
func buildRequest(url string, method string, headers string, body string) *http.Request {
	// 反序列化headers
	tmp := &http.Header{}
	_ = json.Unmarshal([]byte(headers), tmp)

	// Reader化body
	bodyReader := strings.NewReader(body)

	// 构造
	req, _ := http.NewRequest(method, url, bodyReader)
	// 添加Headers
	req.Header = *tmp

	return req
}

func isWebvpn(url string) bool {
	return strings.Contains(url, "webvpn.neu.edu.cn")
}
