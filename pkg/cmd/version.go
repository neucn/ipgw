package cmd

import (
	"github.com/neucn/ipgw"
	"github.com/neucn/ipgw/pkg/console"
	"github.com/urfave/cli/v2"
)

var (
	VersionCommand = &cli.Command{
		Name:  "version",
		Usage: "显示版本及构建信息",
		Action: func(ctx *cli.Context) error {
			console.InfoF("ipgw %s+%s\n", ipgw.Version, ipgw.Build)
			return nil
		},
	}
)
