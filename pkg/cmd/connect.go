package cmd

import (
	"os/exec"

	"github.com/neucn/ipgw/pkg/console"
	"github.com/urfave/cli/v2"
)

var (
	ConnectCommand = &cli.Command{
		Name:      "connect",
		Usage:     "connect to a WLAN with a specific SSID, default NEU",
		ArgsUsage: "[WLAN SSID]",
		Action: func(ctx *cli.Context) error {
			sids := ctx.Args().Slice()
			if len(sids) > 1 {
				console.InfoL("wrong arguments!")
				return nil
			}
			ssid := "NEU"
			if len(sids) == 1 {
				ssid = sids[0]
			}
			console.InfoL("Connecting...")
			out, err := exec.Command("powershell", "netsh wlan connect name="+ssid).Output()
			if err != nil {
				console.InfoF("connect failed. Details: %s\n", err.Error())
			}
			if out != nil {
				console.InfoF(string(out))
			}
			return nil
		},
	}
)
