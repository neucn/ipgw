package cmd

import (
	"errors"
	"fmt"
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/urfave/cli/v2"
)

var (
	LogoutCommand = &cli.Command{
		Name:  "logout",
		Usage: "退出登陆",
		Action: func(ctx *cli.Context) error {
			h := handler.NewIpgwHandler()
			connected, loggedIn := h.IsConnectedAndLoggedIn()
			if !connected {
				return errors.New("未连接到校园网")
			}
			if !loggedIn {
				return errors.New("未登陆")
			}
			info := h.GetInfo()
			if err := h.Logout(); err != nil {
				return fmt.Errorf("无法退出登陆账户 '%s':\n\t%v", info.Username, err)
			}
			console.InfoF("退出登陆 '%s' 成功\n", info.Username)
			return nil
		},
		OnUsageError: onUsageError,
	}
)
