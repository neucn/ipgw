package login

import (
	"ipgw/api/code"
	. "ipgw/api/global"
	"ipgw/ctx"
	. "ipgw/lib"
	"net/http"
	"net/url"
)

func getCookieAfterLoginWithUP(c *ctx.Ctx, vpn bool) {
	// 登陆并过滤登陆失败的情况
	LoginWithUP(c, vpn)
	// 通过筛选，检查cookie
	cookie := getCookie(c, vpn)
	if len(cookie) < 1 {
		Fatal(code.LoginFail)
	}

	// 输出cookie
	InfoLine(cookie)
}

func getCookieAfterLoginWithC(c *ctx.Ctx, vpn bool) {
	// 登陆并过滤登陆失败的情况
	LoginWithC(c, vpn)
	// 如果成功通过了过滤，说明登陆成功
	// 输出Cookie
	InfoLine(c.User.Cookie.Value)
}

func getCookie(c *ctx.Ctx, vpn bool) (cookie string) {
	var cookies []*http.Cookie
	if vpn {
		cookies = c.Client.Jar.Cookies(&url.URL{
			Scheme: "https",
			Host:   "pass-443.webvpn.neu.edu.cn",
			Path:   "/tpass/",
		})
	} else {
		cookies = c.Client.Jar.Cookies(&url.URL{
			Scheme: "https",
			Host:   "pass.neu.edu.cn",
			Path:   "/tpass/",
		})
	}

	// 只需要CASTGC
	for _, v := range cookies {
		if v.Name == "CASTGC" {
			return v.Value
		}
	}
	return ""
}
