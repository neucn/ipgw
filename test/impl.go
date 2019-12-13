package test

import (
	"fmt"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"strings"
	"time"
)

func testImpl() {
	client := ctx.GetClient()
	client.Timeout = time.Second

	if cfg.FullView {
		fmt.Println(testNetTip)
	}

	// 测试是否连接上校园网
	_, err := client.Get("https://ipgw.neu.edu.cn")

	if err != nil {
		if cfg.FullView {
			fmt.Printf("ipgw.neu.edu.cn: %v\n", err)
		}
		if strings.Contains(err.Error(), "no such host") {
			// 没有联网
			fmt.Println(testNoInternet)
			return
		} else if strings.Contains(err.Error(), "Client.Timeout") {
			// 没有连校园网
			fmt.Println(testNotConnect)
			return
		}
		fmt.Printf(testUnexpectedErr, err)
		return
	}

	fmt.Println(testConnected)

	if cfg.FullView {
		fmt.Println(testLoggedTip)
	}

	// 测试是否登陆校园网
	_, err = client.Get("https://baidu.com")

	if err != nil {
		if cfg.FullView {
			fmt.Printf("baidu.com: %v\n", err)
		}
		// todo 还有可能证书错误，再观望观望
		if strings.Contains(err.Error(), "Client.Timeout") {
			// 未登陆
			fmt.Println(testNotLoggedIn)
			return
		}
		fmt.Printf(testUnexpectedErr, err)
		return
	}

	fmt.Println(testLoggedIn)
}
