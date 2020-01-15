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
		FatalL(failGetResp)
	}

	switch title[0][1] {
	case "智慧东大--统一身份认证":
		FatalL(failWrongUOrP)
	case "智慧东大":
		FatalL(failWrongSetting)
	case "系统提示":
		FatalL(failBanned)
	}
}

// 【Cookie登陆】根据title判断是否登陆成功，若不成功则返回false
func LoginStatusFilterC(body string) (ok bool) {
	// 匹配出title
	titleExp := regexp.MustCompile(`<title>(.+?)</title>`)
	title := titleExp.FindAllStringSubmatch(body, -1)
	if len(title) < 1 {
		ErrorL(failGetResp)
		return false
	}

	if title[0][1] == "智慧东大--统一身份认证" {
		ErrorL(failCookieExpired)
		return false
	}

	return true
}
