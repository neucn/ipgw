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
		Usage: "logout ipgw",
		Action: func(ctx *cli.Context) error {
			h := handler.NewIpgwHandler()
			connected, loggedIn := h.IsConnectedAndLoggedIn()
			if !connected {
				return errors.New("不在校园网IP段中")
			}
			if !loggedIn {
				return errors.New("还未登陆")
			}
			info := h.GetInfo()
			if err := h.Logout(); err != nil {
				return fmt.Errorf("无法注销账号 '%s':\n\t%v", info.Username, err)
			}
			console.InfoF("注销账号 '%s' 成功\n", info.Username)
			return nil
		},
		OnUsageError: onUsageError,
	}
)
