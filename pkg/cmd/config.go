package cmd

import (
	"fmt"

	"github.com/neucn/ipgw/pkg/console"
	"github.com/urfave/cli/v2"
)

var (
	ConfigCommand = &cli.Command{
		Name:  "config",
		Usage: "管理配置文件",
		Subcommands: []*cli.Command{
			{
				Name:  "account",
				Usage: "管理配置文件中的账户",
				Subcommands: []*cli.Command{
					configAccountAddCommand,
					configAccountDelCommand,
					configAccountSetCommand,
					configAccountListCommand,
				},
				OnUsageError: onUsageError,
			},
		},
		OnUsageError: onUsageError,
	}

	configAccountAddCommand = &cli.Command{
		Name:  "add",
		Usage: "在配置文件中添加账户",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Aliases:  []string{"u"},
				Required: true,
				Usage:    "学号",
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Required: true,
				Usage:    "网关密码",
			},
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "`secret` for stored account",
			},
			&cli.BoolFlag{
				Name:  "default",
				Usage: "添加默认账户",
			},
		},
		Action: func(ctx *cli.Context) error {
			store, err := getStoreHandler(ctx)
			if err != nil {
				return err
			}
			username := ctx.String("username")
			password := ctx.String("password")
			if err = store.Config.AddAccount(
				username,
				password,
				ctx.String("secret")); err != nil {
				return fmt.Errorf("无法添加账户:\n\t%v", err)
			}

			if ctx.Bool("default") {
				store.Config.SetDefaultAccount(username)
			}
			if err = store.Persist(); err != nil {
				return err
			}
			console.InfoF("'%s' 添加成功\n", username)
			return nil
		},
		OnUsageError: onUsageError,
	}

	configAccountDelCommand = &cli.Command{
		Name:  "del",
		Usage: "从配置文件中删除账户",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Aliases:  []string{"u"},
				Required: true,
				Usage:    "要删除的账户学号",
			},
		},
		Action: func(ctx *cli.Context) error {
			store, err := getStoreHandler(ctx)
			if err != nil {
				return err
			}
			username := ctx.String("username")

			if !store.Config.DelAccount(username) {
				return fmt.Errorf("无法删除账户:\n\t'%s' 账户没有找到", username)
			}

			if err = store.Persist(); err != nil {
				return err
			}
			console.InfoF("'%s' 删除成功\n", username)
			return nil
		},
		OnUsageError: onUsageError,
	}

	configAccountSetCommand = &cli.Command{
		Name:  "set",
		Usage: "编辑配置文件中的账户",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Aliases:  []string{"u"},
				Required: true,
				Usage:    "要编辑的账户学号",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "新的登陆网关密码",
			},
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "new `secret` for stored account, must be used with --password, -p",
			},
			&cli.BoolFlag{
				Name:  "default",
				Usage: "设置为默认账户",
			},
		},
		Action: func(ctx *cli.Context) error {
			store, err := getStoreHandler(ctx)
			if err != nil {
				return err
			}
			username := ctx.String("username")
			account := store.Config.GetAccount(username)
			if account == nil {
				return fmt.Errorf("无法设置账户:\n\t'%s' 账户没有找到", username)
			}

			password := ctx.String("password")
			if password != "" {
				if err = account.SetPassword(ctx.String("password"), []byte(ctx.String("secret"))); err != nil {
					return fmt.Errorf("无法设置账户:\n\t'%s' 账户没有找到", username)
				}
			}

			if ctx.Bool("default") {
				store.Config.SetDefaultAccount(username)
			}

			if err = store.Persist(); err != nil {
				return err
			}
			console.InfoF("'%s' 编辑成功\n", username)
			return nil
		},
		OnUsageError: onUsageError,
	}

	configAccountListCommand = &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "配置文件中的账户列表",
		Action: func(ctx *cli.Context) error {
			store, err := getStoreHandler(ctx)
			if err != nil {
				return err
			}

			for i, account := range store.Config.Accounts {
				console.InfoF("#%d %s", i, account.String())
				if account.Username == store.Config.DefaultAccount {
					console.Info(" - default")
				}
				console.InfoL()
			}

			return nil
		},
		OnUsageError: onUsageError,
	}
)
