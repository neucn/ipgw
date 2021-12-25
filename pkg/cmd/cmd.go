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
	App = &cli.App{
		Name:      "ipgw",
		HelpName:  "ipgw",
		Copyright: "主页:\thttps://github.com/neucn/ipgw\nFeedback:\thttps://github.com/neucn/ipgw/issues/new",
		Commands: []*cli.Command{
			LoginCommand,
			LogoutCommand,
			KickCommand,
			InfoCommand,
			ConfigCommand,
			TestCommand,
			VersionCommand,
			UpdateCommand,
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 0 {
				console.InfoL("未找到命令\n")
				cli.ShowAppHelpAndExit(ctx, 1)
				return nil
			}
			return loginUseDefaultAccount(ctx)
		},
		HideVersion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Hidden:  true,
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"f"},
				Usage:   "载入配置文件",
			},
		},
		OnUsageError: onUsageError,
	}
)

func loginUseDefaultAccount(ctx *cli.Context) error {
	// login use default account
	store, err := getStoreHandler(ctx)
	if err != nil {
		return err
	}
	account := store.Config.GetDefaultAccount()
	if account == nil {
		return errors.New("无储存账户")
	}
	console.InfoF("正在使用账户 '%s'\n", account.Username)
	account.Secret = ctx.String("secret")

	if err = login(handler.NewIpgwHandler(), account); err != nil {
		return fmt.Errorf("登陆失败: \n\t%v", err)
	}
	return nil
}

func getAccountByContext(ctx *cli.Context) (account *model.Account, err error) {
	store, err := getStoreHandler(ctx)
	if err != nil {
		return nil, err
	}

	if c := ctx.String("cookie"); c != "" {
		// use cookie
		account = &model.Account{
			Cookie: c,
		}
	} else if u := ctx.String("username"); u == "" {
		// use stored default account
		if account = store.Config.GetDefaultAccount(); account == nil {
			return nil, errors.New("无储存账户\n\t请提供账号与密码")
		}
		console.InfoF("正在使用账户 '%s'\n", account.Username)
	} else if p := ctx.String("password"); p == "" {
		// use stored account
		if account = store.Config.GetAccount(u); account == nil {
			return nil, fmt.Errorf("账号 '%s' 未找到", u)
		}
	} else {
		// use username and password
		account = &model.Account{
			Username: u,
			Password: p,
		}
	}
	account.Secret = ctx.String("secret")
	return account, nil
}

func getStoreHandler(ctx *cli.Context) (store *handler.StoreHandler, err error) {
	if store, err = handler.NewStoreHandler(ctx.String("config")); err == nil {
		err = store.Load()
	}
	return
}

func onUsageError(ctx *cli.Context, err error, isSubcommand bool) error {
	_, _ = fmt.Fprintf(ctx.App.Writer, "%s\n\n", err.Error())
	if isSubcommand {
		cli.ShowSubcommandHelpAndExit(ctx, 1)
	} else {
		cli.ShowAppHelpAndExit(ctx, 1)
	}
	return nil
}

func init() {
	cli.AppHelpTemplate = `用法:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
{{if .Commands}}
命令:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}
选项:
   {{range .VisibleFlags}}{{.}}
   {{end}}
{{.Copyright}}
`
	cli.CommandHelpTemplate = `{{.Usage}}

用法:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

目录:
   {{.Category}}{{end}}{{if .Description}}

描述:
   {{.Description | nindent 3 | trim}}{{end}}{{if .VisibleFlags}}

选项:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

	cli.SubcommandHelpTemplate = `{{.Usage}}

用法:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Description}}

描述:
   {{.Description | nindent 3 | trim}}{{end}}

命令:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

选项:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
	return
}
