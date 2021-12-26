package cmd

import (
	"errors"
	"fmt"

	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/neucn/ipgw/pkg/model"
	"github.com/urfave/cli/v2"
)

var (
	LoginCommand = &cli.Command{
		Name:  "login",
		Usage: "登陆网关",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "学号（仅在使用非默认账户或者首次储存默认账户时使用）",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "网关登陆密码（仅在账户未储存时需要）",
			},
			&cli.StringFlag{
				Name:    "cookie",
				Aliases: []string{"c"},
				Usage:   "`cookie` item 'session_for%3Asrun_cas_php' from ipgw.neu.edu.cn",
			},
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "`secret` for stored account (required only if secret is not empty)",
			},
			&cli.BoolFlag{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "output account info after login successfully",
			},
		},
		Action: func(ctx *cli.Context) error {
			account, err := getAccountByContext(ctx)
			if err != nil {
				return err
			}
			h := handler.NewIpgwHandler()
			if err = login(h, account); err != nil {
				return fmt.Errorf("登陆失败： \n\t%v", err)
			}
			if ctx.Bool("info") {
				if err = h.FetchUsageInfo(); err != nil {
					return fmt.Errorf("获取信息失败： \n\t%v", err)
				}
				info := h.GetInfo()
				console.InfoF("\tIP\t%16s\n\t余额\t%16s\n\t流量\t%16s\n\t时长\t%16s\n",
					info.IP,
					info.FormattedBalance(),
					info.FormattedTraffic(),
					info.FormattedUsedTime())
			}
			return nil
		},
		OnUsageError: onUsageError,
	}
)

func login(h *handler.IpgwHandler, account *model.Account) error {
	// check logged
	connected, loggedIn := h.IsConnectedAndLoggedIn()
	if !connected {
		return errors.New("未连接到校园网")
	}
	if loggedIn {
		return fmt.Errorf("已登陆为 '%s'", h.GetInfo().Username)
	}
	if err := h.Login(account); err != nil {
		return err
	}
	info := h.GetInfo()
	if info.Username == "" {
		return fmt.Errorf("未知原因")
	}
	console.InfoL("登陆成功")
	return nil
}
