// 存放所有的处理函数

package cas

import (
	. "ipgw/base"
	"regexp"
)

// 【账号登陆】根据title判断是否登陆成功，若不成功则结束
func LoginStatusFilterUP(body string) {
	// 匹配出title
	titleExp := regexp.MustCompile(`<title>(.+?)</title>`)
	title := titleExp.FindAllStringSubmatch(body, -1)
	if len(title) < 1 {
		Fatal(failGetResp)
	}

	switch title[0][1] {
	case "智慧东大--统一身份认证":
		Fatal(failWrongUOrP)
	case "智慧东大":
		Fatal(failWrongSetting)
	case "系统提示":
		Fatal(failBanned)
	}
}

// 【Cookie登陆】根据title判断是否登陆成功，若不成功则返回false
func LoginStatusFilterC(body string) (ok bool) {
	// 匹配出title
	titleExp := regexp.MustCompile(`<title>(.+?)</title>`)
	title := titleExp.FindAllStringSubmatch(body, -1)
	if len(title) < 1 {
		Error(failGetResp)
		return false
	}

	if title[0][1] == "智慧东大--统一身份认证" {
		Error(failCookieExpired)
		return false
	}

	return true
}
