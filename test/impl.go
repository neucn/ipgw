package test

import (
	"fmt"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"ipgw/text"
	"strings"
)

func testImpl() {
	client := ctx.GetClient()

	_, err := client.Get("https://ipgw.neu.edu.cn")

	if cfg.FullView && err != nil {
		fmt.Printf("ipgw.neu.edu.cn: %v\n", err)
	}

	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Println(text.TestNotConnect)
			return
		}
		fmt.Printf(text.TestUnexpectedErr, err)
		return
	}

	_, err = client.Get("https://baidu.com")

	if cfg.FullView && err != nil {
		fmt.Printf("baidu.com: %v\n", err)
	}

	if err != nil {
		if strings.Contains(err.Error(), "x509") {
			fmt.Println(text.TestNotLoggedIn)
			return
		}
		fmt.Printf(text.TestUnexpectedErr, err)
		return
	}

	fmt.Println(text.TestLoggedIn)
}
