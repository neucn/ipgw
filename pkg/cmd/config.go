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
				Usage: "manage accounts stored in config",
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
		Usage: "添加账户至配置文件",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Aliases:  []string{"u"},
				Required: true,
				Usage:    "student number `id`",
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Required: true,
				Usage:    "`password` for pass.neu.edu.cn",
			},
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "`secret` for stored account",
			},
			&cli.BoolFlag{
				Name:  "default",
				Usage: "add the account as the default one",
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
				return fmt.Errorf("无法添加账号:\n\t%v", err)
			}

			if ctx.Bool("default") {
				store.Config.SetDefaultAccount(username)
			}
			if err = store.Persist(); err != nil {
				return err
			}
			console.InfoF("'%s' 添加账号成功\n", username)
			return nil
		},
		OnUsageError: onUsageError,
	}

	configAccountDelCommand = &cli.Command{
		Name:  "del",
		Usage: "从配置文件中删除账号",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Aliases:  []string{"u"},
				Required: true,
				Usage:    "student number `id` to be deleted",
			},
		},
		Action: func(ctx *cli.Context) error {
			store, err := getStoreHandler(ctx)
			if err != nil {
				return err
			}
			username := ctx.String("username")

			if !store.Config.DelAccount(username) {
				return fmt.Errorf("无法删除账号:\n\t'%s' 账号未找到", username)
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
		Usage: "编辑配置文件中的账号",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Aliases:  []string{"u"},
				Required: true,
				Usage:    "student number `id` to be edited",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "new `password` for pass.neu.edu.cn",
			},
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "new `secret` for stored account, must be used with --password, -p",
			},
			&cli.BoolFlag{
				Name:  "default",
				Usage: "set the account as the default one",
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
				return fmt.Errorf("无法设置账号:\n\t'%s' 账号未找到", username)
			}

			password := ctx.String("password")
			if password != "" {
				if err = account.SetPassword(ctx.String("password"), []byte(ctx.String("secret"))); err != nil {
					return fmt.Errorf("无法设置账号:\n\t'%s' 账号未找到", username)
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
		Usage:   "显示在配置文件中的账号",
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
