package list

import (
	"fmt"
	"ipgw/base"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"os"
	"strconv"
	"strings"
)

var CmdList = &base.Command{
	CustomFlags: true,
	UsageLine:   "ipgw list [-v full view] [-s saved] [-u username] [-p password] [-c cookie] [-a all] [-l local info] [-d devices] [-i net info] [-r recharge] [-b bill] [-h history] page",
	Short:       "获取各类信息",
	Long: `提供校园网信息查询功能，默认使用当前登陆的账号的信息
  -s    使用保存的账号查询
  -c    使用cookie查询
  -u    使用指定账号查询，需配合 -p
  -p    使用指定账号查询
  -a    列出所有信息
  -l    列出本地保存的账号及网络信息
  -i    列出校园网套餐信息
  -r	列出充值记录
  -d    列出登陆设备
  -b    列出历史账单
  -h    列出校园网使用日志
  -v    输出详细查询结果及中间信息

  ipgw list
    效果等同于 ipgw list -l
  ipgw list -l
    列出本地保存的账号及会话信息
    包括 已保存账号 Cookie CAS
  ipgw list -a
    效果等同于 ipgw list -birdh 1
    列出当前登陆账号所有信息，必须是使用本工具登陆
  ipgw list -i
    查看当前登陆账号的校园网套餐信息
    包括 套餐 使用流量 使用时长 余额 使用次数
    可使用 -u -p 或 -s 或 -c 查询指定的账号
  ipgw list -r
    列出当前登陆账号的充值记录
    可使用 -u -p 或 -s 或 -c 查询指定的账号
  ipgw list -d
    列出当前登陆账号的已登录设备
    可使用 -u -p 或 -s 或 -c 查询指定的账号
  ipgw list -b
    列出当前登陆账号的历史付费记录
    可使用 -u -p 或 -s 或 -c 查询指定的账号
  ipgw list -h 1
    列出当前登陆账号的使用记录的第一页，每页20条
    可使用 -u -p 或 -s 或 -c 查询指定的账号
  ipgw list -av
    列出所有信息的具体查询结果及中间信息
`,
}

var (
	flags                  = []int32{'s', 'c', 'u', 'p', 'a', 'l', 'i', 'd', 'b', 'h', 'v', 'r'}
	u, p, c                string
	a, i, d, s, h, b, l, r bool
)

func init() {
	CmdList.Flag.BoolVar(&a, "a", false, "列出所有信息")
	CmdList.Flag.BoolVar(&i, "i", false, "列出校园网信息")
	CmdList.Flag.BoolVar(&d, "d", false, "列出登陆设备")
	CmdList.Flag.BoolVar(&l, "l", false, "列出本地保存的信息")
	CmdList.Flag.BoolVar(&b, "b", false, "列出历史付费记录")
	CmdList.Flag.BoolVar(&h, "h", false, "列出校园网使用日志")
	CmdList.Flag.BoolVar(&r, "r", false, "列出校园网充值记录")

	CmdList.Flag.BoolVar(&s, "s", false, "使用保存的账号查询")
	CmdList.Flag.StringVar(&c, "c", "", "使用cookie查询")
	CmdList.Flag.StringVar(&u, "u", "", "使用指定账号查询")
	CmdList.Flag.StringVar(&p, "p", "", "使用指定账号查询")

	CmdList.Flag.BoolVar(&cfg.FullView, "v", false, "输出所有中间信息")

	CmdList.Run = runList // break init cycle
}

func runList(cmd *base.Command, args []string) {
	parse(cmd, args)
	//fmt.Println(a, c, d, h, i, s, cfg.FullView, cmd.Flag.Args())
	x := ctx.GetCtx()
	var ib string
	if u != "" {
		if p == "" {
			fmt.Fprintln(os.Stderr, mustUsePWhenUseU)
			return
		}
		x.User.Username = u
		x.User.Password = p
		ib = fetchIndexBodyByUP(x)
	} else if c != "" {
		x.User.SetCAS(c)
		ib = fetchIndexBodyByCAS(x)
	} else if s {
		x.Load()
		ib = fetchIndexBodyByUP(x)
	} else if a || i || d || h || b || r {
		x.Load()
		ib = fetchIndexBodyByCAS(x)
	}

	if a || i || d || h || b || r {
		processUser(ib)
	} else {
		x.Load()
		processLocal(x)
	}

	if a || i {
		processInfo(ib)
	}
	if a || d {
		processDevice(ib)
	}

	if a || b {
		bb := fetchBillBody()
		processBill(bb)
	}

	if a || r {
		rb := fetchRechargeBody()
		processRecharge(rb)
	}

	if a || h {
		var hb string
		if len(cmd.Flag.Args()) < 1 {
			hb = fetchHistoryBody(1)
		} else {
			p, e := strconv.Atoi(cmd.Flag.Args()[0])
			if e != nil {
				fmt.Fprintln(os.Stderr, wrongPageNotInt)
			}
			hb = fetchHistoryBody(p)
		}

		processHistory(hb)
	}

}

func parse(cmd *base.Command, args []string) {
	separated := make([]string, 0, len(args))
	for _, flagChar := range args {
		if len(flagChar) > 2 && strings.HasPrefix(flagChar, "-") {
		charLoop:
			for _, c := range flagChar[1:] {
				for _, f := range flags {
					if c == f {
						separated = append(separated, "-"+string(c))
						continue charLoop
					}
				}
				fmt.Fprintf(os.Stdout, wrongArgNotFound, string(c))
				cmd.Usage()
			}
			continue
		}
		separated = append(separated, flagChar)
	}
	cmd.Flag.Parse(separated)
}
