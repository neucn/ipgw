package list

import (
	. "ipgw/base"
	"ipgw/ctx"
	"strconv"
	"strings"
)

var CmdList = &Command{
	Name:        "list",
	CustomFlags: true,
	UsageLine:   "ipgw list [-f full] [-v view all] [-s saved] [-u username] [-p password] [-c cookie] [-a all] [-l local info] [-d devices] [-i net info] [-r recharge] [-b bill] [-h history] page",
	Short:       "获取各类信息",
	Long: `提供校园网信息查询功能，默认使用当前登陆的账号的信息
  -s    使用保存的账号查询
  -c    使用cookie查询
  -u    使用指定账号查询，需配合 -p
  -p    使用指定账号查询
  -a    列出所有信息
  -l    列出本地保存的账号及网络信息
  -i    列出校园网套餐信息
  -r    列出充值记录
  -d    列出登陆设备
  -b    列出历史账单
  -h    列出校园网使用日志
  -f    输出所有查询结果的详细信息
  -v    输出所有中间信息

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
  ipgw list -af
    列出所有信息的具体查询结果
  ipgw list -av
    列出中间信息
`,
}

var (
	flags                     = []int32{'s', 'c', 'u', 'p', 'a', 'l', 'i', 'd', 'b', 'h', 'v', 'r', 'f'}
	u, p, c                   string
	a, i, d, s, h, b, l, r, f bool
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

	CmdList.Flag.BoolVar(&f, "f", false, "列出所有信息的具体查询结果")
	CmdList.Flag.BoolVar(&ctx.FullView, "v", false, "输出所有中间信息")

	CmdList.Run = runList
}

func runList(cmd *Command, args []string) {
	// 解析，把缩写形式拆开
	parse(cmd, args)

	x := ctx.NewCtx()

	// Index Body
	var ib string
	if u != "" {
		// 如果使用账号登陆
		if p == "" {
			// 必须要密码
			ErrorL(mustUsePWhenUseU)
			return
		}
		x.User.Username = u
		x.User.Password = p
		ib = fetchIndexBodyByUP(x)
	} else if c != "" {
		// 使用Cookie登陆
		x.User.SetCookie(c)
		ib = fetchIndexBodyByC(x)
	} else if s {
		// 使用保存的账号登陆
		x.Load()
		if x.User.Username == "" {
			FatalL(noStoredAccount)
		}
		ib = fetchIndexBodyByUP(x)
	} else if a || i || d || h || b || r {
		// 若涉及网络请求
		// 使用Cookie登陆
		x.Load()
		ib = fetchIndexBodyByC(x)
	}

	if a || i || d || h || b || r {
		// 若涉及网络请求
		processUser(ib)
	} else {
		// 若不涉及网络请求，直接打印本地信息
		x.Load()
		processLocal(x)
	}

	if a || i {
		// 获取套餐信息
		processInfo(ib)
	}
	if a || d {
		// 获取登陆设备信息
		processDevice(ib)
	}

	if a || b {
		// 获取扣费信息
		bb := fetchBillBody(x)
		processBill(bb)
	}

	if a || r {
		// 获取充值信息
		rb := fetchRechargeBody(x)
		processRecharge(rb)
	}

	if a || h {
		// 获取使用记录
		var hb string
		if len(cmd.Flag.Args()) < 1 {
			// 若没有剩余的参数，即没有指定第几页，默认1
			hb = fetchHistoryBody(x, 1)
		} else {
			p, e := strconv.Atoi(cmd.Flag.Args()[0])
			if e != nil {
				// 检验页数是否为整形
				ErrorL(wrongPageNotInt)
			}
			hb = fetchHistoryBody(x, p)
		}
		processHistory(hb)
	}

}

// 支持flag缩写，算法可能还能优化
func parse(cmd *Command, args []string) {
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
				ErrorF(wrongArgNotFound, string(c))
				cmd.Usage()
			}
			continue
		}
		separated = append(separated, flagChar)
	}
	_ = cmd.Flag.Parse(separated)
}
