package logout

import (
	"fmt"
	"ipgw/base"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"net/http"
	"os"
)

var (
	u, p, c string

	s bool
)

func init() {
	CmdLogout.Flag.StringVar(&u, "u", "", "")
	CmdLogout.Flag.StringVar(&p, "p", "", "")
	CmdLogout.Flag.StringVar(&c, "c", "", "")

	CmdLogout.Flag.BoolVar(&cfg.FullView, "v", false, "")

	CmdLogout.Run = runLogout
}

var CmdLogout = &base.Command{
	UsageLine: "ipgw logout [-u username] [-p password] [-c cookie] [-v full view]",
	Short:     "基础登陆",
	Long: `提供登出校园网关功能
  -u    登出账户
  -p    登出密码
  -c    使用cookie登出
  -v    输出所有中间信息

  ipgw logout -u 学号 -p 密码
    使用指定账号登出网关，必须和当前登陆账号相同
  ipgw logout
    若已经使用-s保存了账户信息，且该账户就是当前登陆的账户，可直接登出
    若没有使用-s保存账户信息，但本次登陆使用的是本工具，也可直接登出
  ipgw logout -c "ST-XXXXXX-XXXXXXXXXXXXXXXXXXXX-tpass"
    使用指定cookie登出
  ipgw logout [arguments] -v
    打印登出过程中的每一步详细信息
`,
}

func init() {
	CmdLogout.Run = runLogout // break init cycle
}

func runLogout(cmd *base.Command, args []string) {
	x := ctx.GetCtx()

	if len(u) > 0 {
		if len(p) == 0 {
			fmt.Fprint(os.Stderr, mustUsePWhenUseU)
			return
		}
		x.User.Username = u
		x.User.Password = p
		logoutWithUP(x)
	} else if len(c) > 0 {
		x.User.Cookie = &http.Cookie{
			Name:   "session_for%3Asrun_cas_php",
			Value:  c,
			Domain: "ipgw.neu.edu.cn",
		}
		ok := logoutWithC(x)
		if !ok {
			os.Exit(2)
		}
	} else {
		x.User.Load(".ipgw")
		// 这就要求不能直接在方法里os.Exit()了
		ok := logoutWithC(x)
		if ok {
			return
		}
		// 若cookie失效，则使用手动
		logoutWithUP(x)
	}
}
