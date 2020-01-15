// 存放对所有业务的拆解函数

package cas

import (
	. "ipgw/base"
	. "ipgw/core/global"
	"ipgw/ctx"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// 登陆 - CAS(根据url, 获取LT, 拼接请求, 校验响应标题), GW(校验是否重复登陆, 退出先前账号, 获取登陆信息)

// 获取lt和postUrl
func GetArgs(c *ctx.Ctx, reqUrl string) (lt, postUrl string) {
	mute := c.Option.Mute
	req, _ := http.NewRequest("GET", reqUrl, nil)
	SendRequestCustomText(c, req, errWhenReadArgs, errNetwork)

	// 读取响应内容
	body := ReadBody(c.Response)

	// 读取lt
	ltExp := regexp.MustCompile(`name="lt" value="(.+?)"`)
	lts := ltExp.FindAllStringSubmatch(body, -1)

	postUrlExp := regexp.MustCompile(`id="loginForm" action="(.+?)"`)
	postUrls := postUrlExp.FindAllStringSubmatch(body, -1)

	if len(lts) < 1 || len(postUrls) < 1 {
		FatalL(failGetArgs)
	}
	lt = lts[0][1]
	postUrl = postUrls[0][1]

	if !mute && ctx.FullView {
		InfoF(successGetArgs, lt)
	}

	return
}

// 登陆
func LoginCAS(c *ctx.Ctx, r *http.Request) {
	// 输出提示
	if !c.Option.Mute && ctx.FullView {
		InfoL(infoSendingLoginRequest)
	}

	// 发送请求
	SendRequest(c, r)
}

// 构造登陆请求
func BuildLoginRequest(c *ctx.Ctx, lt, postUrl, reqUrl string) (req *http.Request) {
	data := "rsa=" + c.User.Username + c.User.Password + lt +
		"&ul=" + strconv.Itoa(len(c.User.Username)) +
		"&pl=" + strconv.Itoa(len(c.User.Password)) +
		"&lt=" + lt +
		"&execution=e1s1" +
		"&_eventId=submit"

	// 构造请求
	req, _ = http.NewRequest("POST",
		GetDomain(reqUrl)+postUrl,
		strings.NewReader(data))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "pass.neu.edu.cn")
	req.Header.Add("Origin", "https://pass.neu.edu.cn")
	req.Header.Add("Referer", reqUrl)
	if c.Option.FakeDevice != "" {
		req.Header.Add("User-Agent", c.Option.FakeDevice)
	}
	return
}

// 伪造UA
func FakeDevice(c *ctx.Ctx, name string) {
	if name == "" {
		return
	}

	var ua string

	// todo 由于移动端会跳转到单独网页，因此暂时把移动端去掉
	switch name {
	case "win", "windows":
		ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36"
	case "linux":
		ua = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"
	case "osx", "darwin":
		ua = "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11"
	//case "ios":
	//	ua = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0_2 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Mobile/15A421 wxwork/2.5.8 MicroMessenger/6.3.22 Language/zh"
	//case "android":
	//	ua = "Mozilla/5.0 (Linux; Android 8.0; DUK-AL20 Build/HUAWEIDUK-AL20; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044353 Mobile Safari/537.36 MicroMessenger/6.7.3.1360(0x26070333) NetType/WIFI Language/zh_CN Process/tools"
	default:
		InfoF(infoDeviceNotFound, name)
	}

	c.Option.FakeDevice = ua
}
