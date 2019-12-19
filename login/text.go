package login

var (
	mustUsePWhenUseU = "使用账号登陆时请务必使用 -p 输入密码"
	noStoredAccount  = "没有已保存的可用账号"

	usingUP        = "使用账号登陆\t%s\n"
	usingC         = "使用Cookie登陆\t%s\n"
	sendingRequest = "发送登陆请求中..."

	errWhenReadLT = "获取参数时错误\t%v\n"
	errNetwork    = "请检查网络连接"

	failGetCookie     = "获取Cookie失败"
	failCookieExpired = "Cookie已失效"
	failLogin         = "登陆失败"
	failGetInfo       = "用户身份获取失败"
	failBalanceOut    = "余额不足 请充值"

	wrongUOrP = "学号或密码错误 请重试"
	wrongLock = "一网通设置有误 请打开网页登陆查看"
	wrongBan  = "一网通系统报错 可能被ban"

	successGetLT       = "获取参数成功\t%s\n"
	successGetCookie   = "获取Cookie成功\t%s\n"
	successGetUsername = "用户身份\t%s\n"
	successLogin       = "登陆成功\t%s\n"
)
