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
				return errors.New("not in campus network")
			}
			if !loggedIn {
				return errors.New("not logged in yet")
			}
			info := h.GetInfo()
			if err := h.Logout(); err != nil {
				return fmt.Errorf("fail to logout account '%s':\n\t%v", info.Username, err)
			}
			console.InfoF("logout account '%s' successfully\n", info.Username)
			return nil
		},
		OnUsageError: onUsageError,
	}
)
