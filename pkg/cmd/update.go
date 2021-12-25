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
		Usage: "检查最新版本",
		Action: func(ctx *cli.Context) error {
			h := handler.NewUpdateHandler()
			newer, err := h.CheckLatestVersion()
			if err != nil {
				return err
			}
			if !newer {
				console.InfoL("已是最新版本")
				return nil
			}
			err = h.Update()
			if err != nil {
				return fmt.Errorf("无法获取更新:\n\t%v", err)
			}
			console.InfoL("更新成功")
			return nil
		},
	}
)
