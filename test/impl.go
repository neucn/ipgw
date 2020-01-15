package test

import (
	. "ipgw/base"
	"ipgw/ctx"
	"net/http"
	"strings"
	"time"
)

func test() {
	// 直接实例化获取client
	client := &http.Client{Timeout: time.Second}

	if ctx.FullView {
		InfoL(testingNet)
	}

	// 测试是否连接上校园网
	_, err := client.Get("https://ipgw.neu.edu.cn")

	if err != nil {
		if ctx.FullView {
			ErrorF("ipgw.neu.edu.cn: %v\n", err)
		}
		if strings.Contains(err.Error(), "no such host") {
			// 没有联网
			InfoL(noInternet)
			return
		} else if strings.Contains(err.Error(), "Client.Timeout") {
			// 没有连校园网
			InfoL(notConnect)
			return
		}
		ErrorF(errUnexpected, err)
		return
	}

	InfoL(connected)

	if ctx.FullView {
		InfoL(testingLogin)
	}

	// 测试是否登陆校园网
	_, err = client.Get("https://baidu.com")

	if err != nil {
		if ctx.FullView {
			ErrorF("baidu.com: %v\n", err)
		}
		// todo 还有可能证书错误，再观望观望
		if strings.Contains(err.Error(), "Client.Timeout") {
			// 未登陆
			InfoL(notLoggedIn)
			return
		}
		ErrorF(errUnexpected, err)
		return
	}

	InfoL(loggedIn)
}
