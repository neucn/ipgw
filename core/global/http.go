package global

import (
	"io/ioutil"
	. "ipgw/base"
	"ipgw/ctx"
	"net/http"
	"regexp"
)

// 发送请求
func SendRequest(c *ctx.Ctx, r *http.Request) {
	SendRequestCustomText(c, r, errRequest, errNetwork)
}

func SendRequestCustomText(c *ctx.Ctx, r *http.Request, errText, fatalText string) {
	resp, err := c.Client.Do(r)
	if err != nil {
		if !c.Option.Mute && ctx.FullView {
			ErrorF(errText, err)
		}
		FatalL(fatalText)
	}
	c.Response = resp
}

// 读取响应体，不可对一个Response对象读多次
func ReadBody(resp *http.Response) (body string) {
	res, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return string(res)
}

// url中提取域名
func GetDomain(url string) (domain string) {
	domainExp := regexp.MustCompile(`(https?://.+?)/`)
	domains := domainExp.FindAllStringSubmatch(url, -1)
	if len(domains) < 1 {
		// todo 有待考量是否提示请求地址有误
		return "https://pass.neu.edu.cn"
	}
	return domains[0][1]
}
