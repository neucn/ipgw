package cmd

import (
	"errors"
	"fmt"
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/neucn/ipgw/pkg/model"
	"github.com/urfave/cli/v2"
)

var (
	LoginCommand = &cli.Command{
		Name:  "login",
		Usage: "login ipgw",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "student number `id` (required only if not use the default or first stored account)",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "`password` for pass.neu.edu.cn (required only if account is not stored)",
			},
			&cli.StringFlag{
				Name:    "cookie",
				Aliases: []string{"c"},
				Usage:   "`cookie` item 'session_for%3Asrun_cas_php' from ipgw.neu.edu.cn",
			},
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "`secret` for stored account (required only if secret is not empty)",
			},
			&cli.BoolFlag{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "output account info after login successfully",
			},
		},
		Action: func(ctx *cli.Context) error {
			h := handler.NewIpgwHandler()
			// check logged
			connected, loggedIn := h.IsConnectedAndLoggedIn()
			if !connected {
				return errors.New("not in campus network")
			}
			if loggedIn {
				return fmt.Errorf("already logged in as '%s'", h.GetInfo().Username)
			}

			store, err := getStoreHandler(ctx)
			if err != nil {
				return err
			}

			var account *model.Account
			if c := ctx.String("cookie"); c != "" {
				// use cookie
				account = &model.Account{
					Cookie: c,
				}
			} else if u := ctx.String("username"); u == "" {
				// use stored default account
				if account = store.Config.GetDefaultAccount(); account == nil {
					return errors.New("no stored account\n\tplease provide username and password")
				}
				console.InfoF("using account '%s'\n", account.Username)
			} else if p := ctx.String("password"); p == "" {
				// use stored account
				if account = store.Config.GetAccount(u); account == nil {
					return fmt.Errorf("account '%s' not found", u)
				}
			} else {
				// use username and password
				account = &model.Account{
					Username: u,
					Password: p,
				}
			}
			account.Secret = ctx.String("secret")

			if err = h.Login(account); err != nil {
				return fmt.Errorf("fail to login:\n\t%v", err)
			}
			info := h.GetInfo()
			if info.Username == "" {
				return fmt.Errorf("fail to login")
			}
			console.InfoL("login successfully")
			if ctx.Bool("info") {
				console.InfoF("\tIP\t%16s\n\t余额\t%16s\n\t流量\t%16s\n\t时长\t%16s\n",
					info.IP,
					info.FormattedBalance(),
					info.FormattedTraffic(),
					info.FormattedUsedTime())
				if info.Overdue {
					console.InfoL("\t状态\t已欠费")
				}
			}
			return nil
		},
		OnUsageError: onUsageError,
	}
)
