package cmd

import (
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/urfave/cli/v2"
)

var (
	KickCommand = &cli.Command{
		Name:                   "kick",
		Usage:                  "通过SID使特定设备下线",
		ArgsUsage:              "[sid list]",
		UseShortOptionHandling: true,
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
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "账户密保问题（仅在未设置时需要）",
			},
		},
		Action: func(ctx *cli.Context) error {
			sids := ctx.Args().Slice()
			if len(sids) == 0 {
				console.InfoL("无sid")
				return nil
			}
			account, err := getAccountByContext(ctx)
			if err != nil {
				return err
			}
			h := handler.NewIpgwHandler()
			password, err := account.GetPassword()
			if err != nil {
				return err
			}
			if err = h.NEUAuth(account.Username, password); err != nil {
				return err
			}
			for _, sid := range sids {
				result, err := h.Kick(sid)
				if result {
					console.InfoF("#%s: 完成\n", sid)
				} else {
					console.InfoF("#%s: 失败\n", sid)
					if err != nil {
						console.InfoF("\t%v\n", err)
					}
				}
			}
			return nil
		},
	}
)
