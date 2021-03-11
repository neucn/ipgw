package cmd

import (
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/urfave/cli/v2"
)

var (
	KickCommand = &cli.Command{
		Name:      "kick",
		Usage:     "logout any specific device by SID",
		ArgsUsage: "[sid list]",
		Action: func(ctx *cli.Context) error {
			sids := ctx.Args().Slice()
			if len(sids) == 0 {
				console.InfoL("no sid")
				return nil
			}
			h := handler.NewIpgwHandler()
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
