package api

import (
	"ipgw/base"
	"ipgw/core/cas"
	. "ipgw/core/global"
	"ipgw/ctx"
	. "ipgw/lib"
	"net/http"
	"net/url"
	"regexp"
)

var Login = &base.Command{}

var (
	u, p, c string

	v bool
)

func init() {
	// 使用账号密码
	Login.Flag.StringVar(&u, "u", "", "")
	Login.Flag.StringVar(&p, "p", "", "")
	// 使用指定Cookie
	Login.Flag.StringVar(&c, "c", "", "")
	// 使用指定UA（由于没有完成移动端网页的适配，因此暂时不支持指定UA
	// 目前来看似乎没有移动端网页上有独有功能的情况
	//Login.Flag.StringVar(&d, "d", "", "")

	// todo 虽然list命令里默认使用Cookie登陆，使用保存的账号需要-s，但是这里是否采用这种做法还有待商榷

	// 使用web vpn登陆
	Login.Flag.BoolVar(&v, "v", false, "")

	Login.Run = runAPILogin
}

func runAPILogin(cmd *base.Command, args []string) {
	// 本命令作用： 对于指定账号或配置而言，获取登陆成功后的Cookie；对于指定Cookie而言，验证Cookie有效性；
	// 对于校园网，Cookie为CASTGC
	// 对于webvpn，Cookie为webvpn_username

	// 本命令命令码 1
	// 本命令具体错误码:
	// 1 - 用户名与密码未同时出现
	// 2 - 用户名或密码错误
	// 3 - Cookie已失效
	// 4 - 登陆失败
	// 5 - 无已保存账号

	// 解析Flag
	_ = cmd.Flag.Parse(args)

	// 新建上下文
	x := ctx.NewCtx()

	if len(u) > 0 {
		// 账号密码
		if len(p) == 0 {
			Fatal(loginNoPassword)
		}
		x.User.Username = u
		x.User.Password = p
		loginWithUP(x, v)
	} else if len(c) > 0 {
		// Cookie
		x.User.SetCookie(c)
		loginWithC(x, v)
	} else {
		// 保存的账号
		LoadUser(x)
		if x.User.Username == "" {
			Fatal(loginNoStoredUser)
		}
		loginWithUP(x, v)
	}

}

func loginWithUP(c *ctx.Ctx, vpn bool) {
	// 获取请求地址
	reqUrl := getReqUrl(vpn)

	// 获取lt和postUrl
	lt, postUrl := getArgs(c, reqUrl)

	// 构造请求
	req := cas.BuildLoginRequest(c, lt, postUrl, reqUrl)

	// 发送请求
	sendRequest(c, req)

	// 获取响应
	body := ReadBody(c.Response)

	// 筛选
	loginStatusFilterUP(body)

	// 通过筛选，检查cookie
	cookie := getCookie(c, vpn)
	if len(cookie) < 1 {
		Fatal(loginFail)
	}

	// 输出cookie
	InfoLine(cookie)
}

func loginWithC(c *ctx.Ctx, vpn bool) {
	// 获取请求地址
	reqUrl := getReqUrl(vpn)

	// 构造请求
	req, _ := http.NewRequest("GET", reqUrl, nil)

	// 使用Cookie
	req.Header.Add("Cookie", "CASTGC="+c.User.Cookie.Value)

	// 发送请求
	sendRequest(c, req)

	// 读取响应体
	body := ReadBody(c.Response)

	// 过滤
	loginStatusFilterC(body)
	// 如果成功通过了过滤，说明登陆成功
	// 输出Cookie
	InfoLine(c.User.Cookie.Value)
}

// 根据是否webvpn获取请求地址
func getReqUrl(vpn bool) (reqUrl string) {
	if vpn {
		reqUrl = "https://pass-443.webvpn.neu.edu.cn/tpass/login?service=https%3A%2F%2Fwebvpn.neu.edu.cn%2Fusers%2Fauth%2Fcas%2Fcallback%3Furl"
	} else {
		reqUrl = "https://pass.neu.edu.cn/tpass/login"
	}
	return
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

// 获取Lt与postUrl
func getArgs(c *ctx.Ctx, reqUrl string) (lt, postUrl string) {
	req, _ := http.NewRequest("GET", reqUrl, nil)
	sendRequest(c, req)
	// 读取响应内容
	body := ReadBody(c.Response)

	// 读取lt
	ltExp := regexp.MustCompile(`name="lt" value="(.+?)"`)
	lts := ltExp.FindAllStringSubmatch(body, -1)

	postUrlExp := regexp.MustCompile(`id="loginForm" action="(.+?)"`)
	postUrls := postUrlExp.FindAllStringSubmatch(body, -1)

	if len(lts) < 1 || len(postUrls) < 1 {
		Fatal(globalNetError)
	}
	lt = lts[0][1]
	postUrl = postUrls[0][1]
	return
}

// 【账号登陆】根据title判断是否登陆成功，若不成功则结束
func loginStatusFilterUP(body string) {
	// 匹配出title
	titleExp := regexp.MustCompile(`<title>(.+?)</title>`)
	title := titleExp.FindAllStringSubmatch(body, -1)
	if len(title) < 1 {
		Fatal(globalNetError)
	}

	switch title[0][1] {
	case "智慧东大--统一身份认证":
		Fatal(loginWrongUP)
	case "智慧东大", "系统提示":
		Fatal(loginBanned)
	}
}

// 【Cookie登陆】根据title判断是否登陆成功，若不成功则返回false
func loginStatusFilterC(body string) {
	// 匹配出title
	titleExp := regexp.MustCompile(`<title>(.+?)</title>`)
	title := titleExp.FindAllStringSubmatch(body, -1)
	if len(title) < 1 {
		Fatal(globalNetError)
	}

	if title[0][1] == "智慧东大--统一身份认证" {
		Fatal(loginCookieExpired)
	}
}
