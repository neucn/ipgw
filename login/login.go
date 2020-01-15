package login

import (
	. "ipgw/base"
	"ipgw/core/cas"
	"ipgw/core/gw"
	"ipgw/ctx"
)

var (
	u, p, c, d string

	s, i bool
)

func init() {
	CmdLogin.Flag.StringVar(&u, "u", "", "")
	CmdLogin.Flag.StringVar(&p, "p", "", "")
	CmdLogin.Flag.StringVar(&c, "c", "", "")
	CmdLogin.Flag.StringVar(&d, "d", "", "")

	CmdLogin.Flag.BoolVar(&s, "s", false, "")
	CmdLogin.Flag.BoolVar(&i, "i", false, "")
	CmdLogin.Flag.BoolVar(&ctx.FullView, "v", false, "")

	CmdLogin.Run = runLogin
}

var CmdLogin = &Command{
	Name:      "login",
	UsageLine: "ipgw login [-u username] [-p password] [-s save] [-c cookie] [-d device] [-i info] [-v view all] ",
	Short:     "基础登陆",
	Long: `提供登陆校园网关功能
  -u    登陆账号
  -p    登陆密码
  -s    保存该账号
  -c    使用cookie登陆
  -d    使用指定设备信息
  -i    登陆后输出账号信息
  -v    输出所有中间信息

  ipgw
    效果等同于 ipgw login -i
    [推荐] 在已经使用-s保存了账号信息的情况下，直接执行ipgw即可完成登陆
  ipgw login -u 学号 -p 密码
    使用指定账号登陆网关
  ipgw login -u 学号 -p 密码 -s
    本次登陆的账号信息将被保存在用户目录下的.ipgw文件中
  ipgw login
    在已经使用-s保存了账号信息的情况下，可以直接使用已经保存的账号登录
  ipgw login -c "ST-XXXXXX-XXXXXXXXXXXXXXXXXXXX-tpass"
    使用指定cookie登陆
  ipgw login -d linux
    使用指定设备信息登陆，可选的有win linux osx，默认使用匿名设备信息
  ipgw login -i
    登陆成功后输出账号信息，包括账号余额、已使用时长、已使用流量等
  ipgw login [arguments] -v
    打印登陆过程中的每一步信息
`,
}

func runLogin(cmd *Command, args []string) {
	// 获取上下文对象
	x := ctx.NewCtx()

	// 处理设备信息
	cas.FakeDevice(x, d)

	if len(u) > 0 {
		if len(p) == 0 {
			FatalL(mustUsePWhenUseU)
		}
		x.User.Username = u
		x.User.Password = p
		loginWithUP(x)
	} else if len(c) > 0 {
		// 暂定为使用网关的cookie，这样避免了一网通的转发和失效时间
		x.Net.SetCookie(c)
		loginWithC(x)
	} else {
		// 载入配置文件
		x.Load()
		// 没有保存的账号
		if x.User.Username == "" {
			FatalL(noStoredAccount)
		}
		loginWithUP(x)
	}

	// 当直接ipgw进来的args是nil
	if i || args == nil {
		// 获取账号信息
		gw.GetNetInfo(x)
		// 输出
		x.Net.Print()
	}

	if s {
		// 若s，把账号密码Cookie写进文件里
		x.SaveAll()
	} else {
		// 否则只写入Session
		x.SaveSession()
	}
}
