package cmd

import (
	"fmt"

	"github.com/neucn/ipgw/pkg/console"
	"github.com/urfave/cli/v2"
)

var (
	ConfigCommand = &cli.Command{
		Name:  "config",
		Usage: "manage config",
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
		Usage: "add account into config",
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
				return fmt.Errorf("fail to add account:\n\t%v", err)
			}

			if ctx.Bool("default") {
				store.Config.SetDefaultAccount(username)
			}
			if err = store.Persist(); err != nil {
				return err
			}
			console.InfoF("'%s' added successfully\n", username)
			return nil
		},
		OnUsageError: onUsageError,
	}

	configAccountDelCommand = &cli.Command{
		Name:  "del",
		Usage: "delete account from config",
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
				return fmt.Errorf("fail to delete account:\n\t'%s' not found", username)
			}

			if err = store.Persist(); err != nil {
				return err
			}
			console.InfoF("'%s' deleted successfully\n", username)
			return nil
		},
		OnUsageError: onUsageError,
	}

	configAccountSetCommand = &cli.Command{
		Name:  "set",
		Usage: "edit account in config",
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
				return fmt.Errorf("fail to set account:\n\t'%s' not found", username)
			}

			password := ctx.String("password")
			if password != "" {
				if err = account.SetPassword(ctx.String("password"), []byte(ctx.String("secret"))); err != nil {
					return fmt.Errorf("fail to set password:\n\t'%v'", err)
				}
			}

			if ctx.Bool("default") {
				store.Config.SetDefaultAccount(username)
			}

			if err = store.Persist(); err != nil {
				return err
			}
			console.InfoF("'%s' edited successfully\n", username)
			return nil
		},
		OnUsageError: onUsageError,
	}

	configAccountListCommand = &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "list accounts in config",
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
