package global

import (
	"encoding/base64"
	"io/ioutil"
	"ipgw/api/code"
	. "ipgw/base"
	"ipgw/core/cas"
	"ipgw/core/global"
	"ipgw/ctx"
	. "ipgw/lib"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// 发送请求
func SendRequest(c *ctx.Ctx, r *http.Request) {
	resp, err := c.Client.Do(r)
	if err != nil {
		Fatal(code.GlobalNetError)
	}
	c.Response = resp
}

// 载入一网通账号和密码
func LoadUser(c *ctx.Ctx) {
	// 准备读取
	path, err := GetConfigPath(SavePath)
	if err != nil {
		Fatal(code.GlobalFailLoad)
	}

	// 读取
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		Fatal(code.GlobalFailLoad)
	}
	content := string(bytes)

	// 分割
	lines := strings.Split(content, LineDelimiter)
	if len(lines) < 2 {
		Fatal(code.GlobalFailLoad)
	}

	// 载入用户信息部分
	user := strings.Split(lines[0], PartDelimiter)
	if len(user) < 3 {
		Fatal(code.GlobalFailLoad)
	}

	// [b64(username), b64(password), CAS Cookie]
	username, err := base64.StdEncoding.DecodeString(user[0])
	c.User.Username = string(username)

	password, err := base64.StdEncoding.DecodeString(user[1])
	c.User.Password = string(password)
}

// 【账号登陆】根据title判断是否登陆成功，若不成功则结束并报错
func loginStatusFilterUP(body string) {
	// 匹配出title
	titleExp := regexp.MustCompile(`<title>(.+?)</title>`)
	title := titleExp.FindAllStringSubmatch(body, -1)

	if len(title) < 1 {
		return
	}

	switch title[0][1] {
	case "智慧东大--统一身份认证":
		Fatal(code.LoginWrongUP)
	case "智慧东大", "系统提示":
		Fatal(code.LoginBanned)
	}
}

// 【Cookie登陆】根据title判断是否登陆成功，若不成功则结束并报错
func loginStatusFilterC(body string) {
	// 匹配出title
	titleExp := regexp.MustCompile(`<title>(.+?)</title>`)
	title := titleExp.FindAllStringSubmatch(body, -1)

	if len(title) > 0 && title[0][1] == "智慧东大--统一身份认证" {
		Fatal(code.LoginCookieExpired)
	}
}

// 获取Lt与postUrl，获取失败则结束并报错
func getArgs(c *ctx.Ctx, reqUrl string) (lt, postUrl string) {
	req, _ := http.NewRequest("GET", reqUrl, nil)
	SendRequest(c, req)
	// 读取响应内容
	body := global.ReadBody(c.Response)

	// 读取lt
	ltExp := regexp.MustCompile(`name="lt" value="(.+?)"`)
	lts := ltExp.FindAllStringSubmatch(body, -1)

	postUrlExp := regexp.MustCompile(`id="loginForm" action="(.+?)"`)
	postUrls := postUrlExp.FindAllStringSubmatch(body, -1)

	if len(lts) < 1 || len(postUrls) < 1 {
		Fatal(code.GlobalNetError)
	}
	lt = lts[0][1]
	postUrl = postUrls[0][1]
	return
}

// 根据是否webvpn获取请求地址
func getReqUrl(serviceUrl string, vpn bool) (reqUrl string) {
	if vpn {
		return "https://pass-443.webvpn.neu.edu.cn/tpass/login?service=https%3A%2F%2Fwebvpn.neu.edu.cn%2Fusers%2Fauth%2Fcas%2Fcallback%3Furl"
	} else {
		reqUrl = "https://pass.neu.edu.cn/tpass/login?service="
		if len(serviceUrl) > 0 {
			reqUrl += url.QueryEscape(serviceUrl)
		} else {
			reqUrl += "https%3A%2F%2Fportal.neu.edu.cn%2Ftp_up%2F"
		}
		return
	}
}

// 使用Cookie登陆，若失败则直接结束程序并输出错误码，因此不需要返回是否登陆成功
func LoginWithC(c *ctx.Ctx, serviceUrl string, vpn bool) {
	// 获取请求地址
	reqUrl := getReqUrl(serviceUrl, vpn)

	// 构造请求
	req, _ := http.NewRequest("GET", reqUrl, nil)

	// 使用Cookie
	req.Header.Add("Cookie", "CASTGC="+c.User.Cookie.Value)

	// 发送请求
	SendRequest(c, req)

	// 读取响应体
	body := global.ReadBody(c.Response)

	// 过滤
	loginStatusFilterC(body)
}

// 使用账号登陆，若失败则直接结束程序并输出错误码，因此不需要返回是否登陆成功
func LoginWithUP(c *ctx.Ctx, serviceUrl string, vpn bool) {
	// 获取请求地址
	reqUrl := getReqUrl(serviceUrl, vpn)

	// 获取lt和postUrl
	lt, postUrl := getArgs(c, reqUrl)

	// 构造请求
	req := cas.BuildLoginRequest(c, lt, postUrl, reqUrl)

	// 发送请求
	SendRequest(c, req)

	// 获取响应
	body := global.ReadBody(c.Response)

	// 筛选
	loginStatusFilterUP(body)
}
