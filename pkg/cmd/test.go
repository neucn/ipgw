package cmd

import (
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/urfave/cli/v2"
)

var (
	TestCommand = &cli.Command{
		Name:  "test",
		Usage: "测试是否连接校园网以及登陆网关",
		Action: func(ctx *cli.Context) error {
			h := handler.NewIpgwHandler()
			connected, loggedIn := h.IsConnectedAndLoggedIn()
			console.Info("校园网:   ")
			if connected {
				console.InfoL("已连接")
			} else {
				console.InfoL("未连接")
			}
			console.Info("网关已登陆:   ")
			if loggedIn {
				console.InfoL("是")
			} else {
				console.InfoL("否")
			}
			return nil
		},
		OnUsageError: onUsageError,
	}
)
