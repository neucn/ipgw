package login

import (
	"ipgw/api/code"
	. "ipgw/api/global"
	. "ipgw/base"
	"ipgw/ctx"
)

var Login = &Command{}

var (
	u, p, c string

	v bool
)

func init() {
	// 使用账号密码
	Login.Flag.StringVar(&u, "u", "", "")
	Login.Flag.StringVar(&p, "p", "", "")
	// 使用指定Cookie
	Login.Flag.StringVar(&c, "c", "", "")
	// 使用指定UA（由于没有完成移动端网页的适配，因此暂时不支持指定UA
	// 目前来看似乎没有移动端网页上有独有功能的情况
	//Login.Flag.StringVar(&d, "d", "", "")

	// todo 虽然list命令里默认使用Cookie登陆，使用保存的账号需要-s，但是这里是否采用这种做法还有待商榷

	// 使用web vpn登陆
	Login.Flag.BoolVar(&v, "v", false, "")

	Login.Run = runAPILogin
}

func runAPILogin(cmd *Command, args []string) {
	// 本命令作用： 对于指定账号或配置而言，获取登陆成功后的Cookie；对于指定Cookie而言，验证Cookie有效性；
	// Cookie均为CASTGC

	// 本命令命名空间 1
	// 本命令具体错误码:
	// 1 - 用户名与密码未同时出现
	// 2 - 无已保存账号
	// 3 - 用户名或密码错误
	// 4 - Cookie已失效
	// 5 - 账号被Ban
	// 6 - 登陆失败

	// 解析Flag
	_ = cmd.Flag.Parse(args)

	// 新建上下文
	x := ctx.NewCtx()

	if len(u) > 0 {
		// 账号密码
		if len(p) == 0 {
			Fatal(code.LoginNoPassword)
		}
		x.User.Username = u
		x.User.Password = p
		getCookieAfterLoginWithUP(x, v)
	} else if len(c) > 0 {
		// Cookie
		x.User.SetCookie(c)
		getCookieAfterLoginWithC(x, v)
	} else {
		// 保存的账号
		LoadUser(x)
		if x.User.Username == "" {
			Fatal(code.LoginNoStoredUser)
		}
		getCookieAfterLoginWithUP(x, v)
	}
}
