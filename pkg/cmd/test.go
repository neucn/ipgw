package cmd

import (
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/urfave/cli/v2"
)

var (
	TestCommand = &cli.Command{
		Name:  "test",
		Usage: "测试当前是否处于校园网环境以及是否已经登陆网关",
		Action: func(ctx *cli.Context) error {
			h := handler.NewIpgwHandler()
			connected, loggedIn := h.IsConnectedAndLoggedIn()
			console.Info("校园网:   ")
			if connected {
				console.InfoL("连接成功")
			} else {
				console.InfoL("断开连接")
			}
			console.Info("ipgw logged in:   ")
			if loggedIn {
				console.InfoL("yes")
			} else {
				console.InfoL("no")
			}
			return nil
		},
		OnUsageError: onUsageError,
	}
)
