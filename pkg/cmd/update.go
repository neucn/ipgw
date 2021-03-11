package cmd

import (
	"fmt"
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/urfave/cli/v2"
)

var (
	UpdateCommand = &cli.Command{
		Name:  "update",
		Usage: "check latest version of ipgw and update",
		Action: func(ctx *cli.Context) error {
			h := handler.NewUpdateHandler()
			newer, err := h.CheckLatestVersion()
			if err != nil {
				return err
			}
			if !newer {
				console.InfoL("already the latest version")
				return nil
			}
			err = h.Update()
			if err != nil {
				return fmt.Errorf("fail to update:\n\t%v", err)
			}
			console.InfoL("update successfully")
			return nil
		},
	}
)
