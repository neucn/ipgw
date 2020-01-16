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

type proxyConfig struct {
	LaunchUrl  string
	ServiceUrl string
	Method     string
	Headers    string
	Body       string
}

// 使用账号密码代理请求
func proxyWithUP(c *ctx.Ctx, config *proxyConfig) {
	// 先登录
	global.LoginWithUP(c, config.ServiceUrl, isWebvpn(config.ServiceUrl))

	// 代理请求
	proxy(c, config)
}

func proxyWithC(c *ctx.Ctx, config *proxyConfig) {
	// 先登录
	global.LoginWithC(c, config.ServiceUrl, isWebvpn(config.ServiceUrl))

	// 代理请求
	proxy(c, config)
}

func proxy(c *ctx.Ctx, config *proxyConfig) {
	// 登陆成功，获取平台Cookie
	if len(config.LaunchUrl) > 0 {
		_, _ = c.Client.Get(config.LaunchUrl)
	}

	// 构造服务请求
	req := buildRequest(config)

	// 访问
	global.SendRequest(c, req)

	Info(ReadBody(c.Response))
}

// 构造请求
func buildRequest(config *proxyConfig) *http.Request {
	// Reader化body
	bodyReader := strings.NewReader(config.Body)

	// 构造
	req, _ := http.NewRequest(config.Method, config.ServiceUrl, bodyReader)

	// 反序列化headers
	var tmp map[string][]string
	_ = json.Unmarshal([]byte(config.Headers), &tmp)
	// 添加Headers
	for k, v := range tmp {
		// 暂时认为header的值数组是没有必要的
		req.Header.Set(k, v[0])
	}

	return req
}

func isWebvpn(url string) bool {
	return strings.Contains(url, "webvpn.neu.edu.cn")
}
