package login

import (
	. "ipgw/base"
	"ipgw/core/cas"
	. "ipgw/core/global"
	"ipgw/core/gw"
	"ipgw/ctx"
	"net/http"
	"net/url"
	"os"
)

func loginWithUP(c *ctx.Ctx) {
	InfoF(usingUP, c.User.Username)

	reqUrl := "https://pass.neu.edu.cn/tpass/login?service=https%3A%2F%2Fipgw.neu.edu.cn%2Fsrun_cas.php%3Fac_id%3D1"
	// 获取必要参数
	lt, postUrl := cas.GetArgs(c, reqUrl)

	// 产生请求
	req := cas.BuildLoginRequest(c, lt, postUrl, reqUrl)

	// 发送请求
	cas.LoginCAS(c, req)

	// 读取响应体
	body := ReadBody(c.Response)

	// 判断登陆状态
	cas.LoginStatusFilterUP(body)

	// 判断是否重复登陆
	rid, rsid := gw.IsLoginRepeatedly(body)
	// 重复登陆
	if len(rid) > 0 {
		gw.Replace(c, rid, rsid)
		body = ReadBody(c.Response)
	}

	// 判断是否欠费
	overdue := gw.IsOverdue(body)
	if overdue {
		FatalL(failOverdue)
	}

	// 获取SID和IP
	// 获取失败即登陆失败
	ok := getSIDAndIP(c, body)
	if !ok {
		FatalL(failLogin)
	}
	// 提取Cookie
	c.ExtractNetCookie()
	c.ExtractUserCookie()

	// 登陆成功
	InfoF(successLogin, c.User.Username)
}

func loginWithC(c *ctx.Ctx) {
	InfoF(usingC, c.User.Cookie.Value)

	// 使用Cookie
	client := c.Client
	client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "ipgw.neu.edu.cn",
	}, []*http.Cookie{c.Net.Cookie})

	if ctx.FullView {
		InfoL(sendingRequest)
	}

	// 构造请求
	req, _ := http.NewRequest("GET", "https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1", nil)

	// 加上伪造UA
	if c.Option.FakeDevice != "" {
		req.Header.Add("User-Agent", c.Option.FakeDevice)
	}

	// 发送请求
	SendRequest(c, req)

	// 读取响应内容
	body := ReadBody(c.Response)

	// 检查标题
	ok := cas.LoginStatusFilterC(body)
	// 未通过
	if !ok {
		os.Exit(2)
	}

	// 判断是否重复登陆
	rid, rsid := gw.IsLoginRepeatedly(body)
	// 重复登陆
	if len(rid) > 0 {
		gw.Replace(c, rid, rsid)
		body = ReadBody(c.Response)
	}

	// 读取学号
	getID(c, body)

	// 判断是否欠费
	overdue := gw.IsOverdue(body)
	if overdue {
		FatalL(failOverdue)
	}

	// 读取SID
	// 获取失败即登录失败
	ok = getSIDAndIP(c, body)
	if !ok {
		FatalL(failLogin)
	}

	// 提取Cookie
	c.ExtractNetCookie()
	c.ExtractUserCookie()

	InfoF(successLogin, c.User.Username)
}

// 对GetSIDAndIP的包装，方便将SID与IP写入ctx
func getSIDAndIP(c *ctx.Ctx, body string) (ok bool) {
	sid, ip := gw.GetSIDAndIP(body)
	// 判断IP是否存在
	if len(ip) < 1 {
		return false
	}
	c.Net.IP = ip
	if ctx.FullView {
		InfoF(successGetIP, ip)
	}

	// 判断SID是否存在
	if len(sid) < 1 {
		return false
	}
	c.Net.SID = sid
	if ctx.FullView {
		InfoF(successGetSID, sid)
	}

	return true
}

// 对GetID的包装，方便将ID写入ctx
func getID(c *ctx.Ctx, body string) {
	id := gw.GetID(body)
	if len(id) < 1 {
		ErrorL(failGetInfo)
	} else {
		c.User.Username = id
		if ctx.FullView {
			InfoF(successGetUsername, c.User.Username)
		}
	}
}
