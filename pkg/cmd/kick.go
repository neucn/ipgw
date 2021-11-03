package cmd

import (
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/urfave/cli/v2"
)

var (
	KickCommand = &cli.Command{
		Name:                   "kick",
		Usage:                  "logout any specific device by SID",
		ArgsUsage:              "[sid list]",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "student number `id` (required only if not use the default or first stored account)",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "`password` for pass.neu.edu.cn (required only if account is not stored)",
			},
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "`secret` for stored account (required only if secret is not empty)",
			},
		},
		Action: func(ctx *cli.Context) error {
			sids := ctx.Args().Slice()
			if len(sids) == 0 {
				console.InfoL("no sid")
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
					console.InfoF("#%s: done\n", sid)
				} else {
					console.InfoF("#%s: fail\n", sid)
					if err != nil {
						console.InfoF("\t%v\n", err)
					}
				}
			}
			return nil
		},
	}
)
