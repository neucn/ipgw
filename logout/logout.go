package logout

import (
	. "ipgw/base"
	"ipgw/ctx"
	"os"
)

var (
	u, p, c string
)

func init() {
	CmdLogout.Flag.StringVar(&u, "u", "", "")
	CmdLogout.Flag.StringVar(&p, "p", "", "")
	CmdLogout.Flag.StringVar(&c, "c", "", "")

	CmdLogout.Flag.BoolVar(&ctx.FullView, "v", false, "")

	CmdLogout.Run = runLogout
}

var CmdLogout = &Command{
	Name:      "logout",
	UsageLine: "ipgw logout [-u username] [-p password] [-c cookie] [-v view all]",
	Short:     "基础登出",
	Long: `提供登出校园网关功能
  -u    登出账号
  -p    登出密码
  -c    使用cookie登出
  -v    输出所有中间信息

  ipgw logout
    若本次登陆是通过本工具，则直接登出
    若直接登出失败，且有未失效的Cookie，将使用Cookie登出
    若Cookie登出失败，且已使用-s保存了账号信息，将使用该账号登出
  ipgw logout -u 学号 -p 密码
    使用指定账号登出网关
  ipgw logout -c "ST-XXXXXX-XXXXXXXXXXXXXXXXXXXX-tpass"
    使用指定cookie登出
  ipgw logout [arguments] -v
    打印登出过程中的每一步详细信息
`,
}

func runLogout(cmd *Command, args []string) {
	x := ctx.NewCtx()

	if len(u) > 0 {
		if len(p) == 0 {
			FatalL(mustUsePWhenUseU)
		}
		x.User.Username = u
		x.User.Password = p
		logoutWithUP(x)
	} else if len(c) > 0 {
		x.User.SetCookie(c)
		// 为了兼容无参数登出，必须使Cookie登出返回一个是否成功的bool然后在这判断并结束程序
		ok := logoutWithC(x)
		if !ok {
			os.Exit(2)
		}
	} else {
		x.Load()
		ok := logoutWithSID(x)
		if ok {
			return
		}

		ok = logoutWithC(x)
		if ok {
			return
		}

		if x.User.Username == "" {
			FatalL(noStoredAccount)
		}
		// 若cookie失效，则使用账号密码
		logoutWithUP(x)
	}
}
