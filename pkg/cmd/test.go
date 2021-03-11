package cmd

import (
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/urfave/cli/v2"
)

var (
	TestCommand = &cli.Command{
		Name:  "test",
		Usage: "test whether is connected to the campus network and whether has logged in ipgw",
		Action: func(ctx *cli.Context) error {
			h := handler.NewIpgwHandler()
			console.Info("campus network:   ")
			if h.IsConnected() {
				console.InfoL("connected")
			} else {
				console.InfoL("disconnected")
			}
			console.Info("ipgw logged in:   ")
			if h.IsConnected() {
				console.InfoL("yes")
			} else {
				console.InfoL("no")
			}
			return nil
		},
		OnUsageError: onUsageError,
	}
)
