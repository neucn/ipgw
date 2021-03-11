package cmd

import (
	"errors"
	"fmt"
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/neucn/ipgw/pkg/model"
	"github.com/urfave/cli/v2"
	"strings"
)

var (
	InfoCommand = &cli.Command{
		Name:                   "info",
		Usage:                  "list account info",
		UseShortOptionHandling: true,
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
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "`secret` for stored account (required only if secret is not empty)",
			},
			&cli.BoolFlag{Name: "all", Aliases: []string{"a"}, Usage: "list all kind of info, equivalent to -lbird"},

			&cli.BoolFlag{Name: "package", Aliases: []string{"i"}, Usage: "print campus network package info"},
			&cli.BoolFlag{Name: "device", Aliases: []string{"d"}, Usage: "print logged-in devices"},
			&cli.IntFlag{Name: "recharge", Aliases: []string{"r"}, Value: 1, Usage: "print the specific `page` of recharge records"},
			&cli.IntFlag{Name: "bill", Aliases: []string{"b"}, Value: 1, Usage: "print the specific `page` of bills"},
			&cli.IntFlag{Name: "log", Aliases: []string{"l"}, Value: 1, Usage: "print the specific `page` of usage logs"},
		},
		Action: func(ctx *cli.Context) error {
			store, err := getStoreHandler(ctx)
			if err != nil {
				return err
			}
			var account *model.Account
			if u := ctx.String("username"); u == "" {
				// use stored default account
				if account = store.Config.GetDefaultAccount(); account == nil {
					return errors.New("no stored account\n\tplease provide username and password")
				}
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
			h := handler.NewDashboardHandler()
			if err := h.Login(account); err != nil {
				return fmt.Errorf("fail to login:\n\t%v", err)
			}
			processInfoPrint(ctx, &infoPrinter{h})
			return nil
		},
		OnUsageError: onUsageError,
	}
)

func processInfoPrint(ctx *cli.Context, printer *infoPrinter) {
	var i, d, r, b, l bool
	if ctx.Bool("all") {
		i = true
		d = true
		r = true
		b = true
		l = true
	} else {
		i = ctx.Bool("package")
		d = ctx.Bool("device")
		r = ctx.IsSet("recharge")
		b = ctx.IsSet("bill")
		l = ctx.IsSet("log")
	}
	printer.PrintBasic()
	if i {
		printer.PrintPackage()
	}
	if d {
		printer.PrintDevices()
	}
	if l {
		printer.PrintLog(ctx.Int("log"))
	}
	if b {
		printer.PrintBills(ctx.Int("bill"))
	}
	if r {
		printer.PrintRecharges(ctx.Int("recharge"))
	}
}

type infoPrinter struct {
	*handler.DashboardHandler
}

func (i *infoPrinter) PrintBasic() {
	defer console.InfoL()
	console.InfoL("# 基本信息")
	b, err := i.GetBasic()
	if err != nil {
		console.InfoL("\t获取失败")
		return
	}
	console.InfoF("\t姓名\t%s\n\t学号\t%s\n", b.Name, b.ID)
}

func (i *infoPrinter) PrintPackage() {
	defer console.InfoL()
	console.InfoL("# 套餐信息")
	pkg, err := i.GetPackage()
	if err != nil {
		console.InfoL("\t获取失败")
		return
	}
	status := make([]string, 0, 2)
	if pkg.Overdue {
		status = append(status, "已欠费")
	}
	if pkg.ExcessPackageTraffic {
		status = append(status, "流量超额")
	}
	if len(status) == 0 {
		status = append(status, "正常")
	}
	console.InfoF("\t套餐\t%sG / %sR\n\t已用\t%s\n\t时长\t%s\n\t次数\t%s\n\t余额\t%sR\n\t状态\t%s\n",
		pkg.PackageTraffic, pkg.PackageCost, pkg.UsedTraffic, pkg.UsedDuration, pkg.UsedTimes, pkg.Balance, strings.Join(status, " "))
}

func (i *infoPrinter) PrintDevices() {
	defer console.InfoL()
	console.InfoL("# 在线设备")
	devices, err := i.GetDevice()
	if err != nil {
		console.InfoL("\t获取失败")
		return
	}
	for _, device := range devices {
		console.InfoF("\tNo.%d\t%s\t%s\t%s\n", device.ID, device.IP, device.StartTime, device.SID)
	}
}

func (i *infoPrinter) PrintRecharges(page int) {
	defer console.InfoL()
	console.InfoL("# 充值记录")
	records, err := i.GetRecharge(page)
	if err != nil {
		console.InfoL("\t获取失败")
		return
	}
	for _, record := range records {
		console.InfoF("\t#%8s\t%s\t%sR\t%s\n", record.ID, record.Time, record.Cost, record.Method)
	}
}

func (i *infoPrinter) PrintLog(page int) {
	defer console.InfoL()
	console.InfoL("# 使用历史")
	records, err := i.GetUsageRecords(page)
	if err != nil {
		console.InfoL("\t获取失败")
		return
	}
	for _, record := range records {
		console.InfoF("\t%s - %s\t%s\t%s\n", record.StartTime, record.EndTime, record.IP, record.Traffic)
	}
}

func (i *infoPrinter) PrintBills(page int) {
	defer console.InfoL()
	console.InfoL("# 扣费记录")
	records, err := i.GetBill(page)
	if err != nil {
		console.InfoL("\t获取失败")
		return
	}
	for _, record := range records {
		console.InfoF("\t#%8s\t%s\t%sR\t%s\n", record.ID, record.Date, record.Cost, record.Traffic)
	}
}
