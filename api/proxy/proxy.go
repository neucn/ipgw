package proxy

import (
	"ipgw/api/code"
	. "ipgw/api/global"
	. "ipgw/base"
	"ipgw/ctx"
)

var Proxy = &Command{}

var (
	u, p, c, l, s, h, b, m string
)

func init() {
	// 使用账号密码
	Proxy.Flag.StringVar(&u, "u", "", "")
	Proxy.Flag.StringVar(&p, "p", "", "")
	// 使用指定Cookie
	Proxy.Flag.StringVar(&c, "c", "", "")
	// 使用指定UA（由于没有完成移动端网页的适配，因此暂时不支持指定UA
	// 目前来看似乎没有移动端网页上有独有功能的情况
	//Proxy.Flag.StringVar(&d, "d", "", "")
	Proxy.Flag.StringVar(&l, "l", "", "")
	Proxy.Flag.StringVar(&s, "s", "", "")
	Proxy.Flag.StringVar(&m, "m", "GET", "")
	Proxy.Flag.StringVar(&h, "h", "{}", "")
	Proxy.Flag.StringVar(&b, "b", "", "")

	Proxy.Run = runAPIProxy
}

func runAPIProxy(cmd *Command, args []string) {
	// 本命令作用： 代为登陆并发送指定的请求，自动判断是否使用webvpn，并返回请求结果（Body

	// 本命令命名空间 2
	// 本命令具体错误码:
	// 1 - 未指定service url

	_ = cmd.Flag.Parse(args)

	// 检查是否填写url
	if len(s) < 1 {
		Fatal(code.ProxyNoServiceUrl)
	}

	// 新建上下文
	x := ctx.NewCtx()

	// 新建proxyConfig
	config := &proxyConfig{
		LaunchUrl:  l,
		ServiceUrl: s,
		Method:     m,
		Headers:    h,
		Body:       b,
	}

	if len(u) > 0 {
		// 账号密码
		if len(p) == 0 {
			Fatal(code.LoginNoPassword)
		}
		x.User.Username = u
		x.User.Password = p
		proxyWithUP(x, config)
	} else if len(c) > 0 {
		// Cookie
		x.User.SetCookie(c)
		proxyWithC(x, config)
	} else {
		// 保存的账号
		LoadUser(x)
		if x.User.Username == "" {
			Fatal(code.LoginNoStoredUser)
		}
		proxyWithUP(x, config)
	}

}
