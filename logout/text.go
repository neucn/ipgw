package logout

var (
	mustUsePWhenUseU = "使用账号登出时请务必使用 -p 输入密码"
	noStoredAccount  = "没有已保存的可用账号"
	balanceOut       = "余额不足"

	errWhenReadLT = "获取参数时错误\t%v\n"
	errNetwork    = "请检查网络连接"

	usingUP        = "使用账号登出\t%s\n"
	usingC         = "使用Cookie登出\t%s\n"
	usingSID       = "使用SID登出\t%s\n"
	sendingRequest = "发送登出请求中..."

	beginLogout = "开始登出\t%s\n"

	successGetLT       = "获取参数成功\t%s\n"
	successGetID       = "成功获得学号\t%s\n"
	successLogout      = "登出成功\t%s\n"
	successLogoutBySID = "登出成功"

	wrongUOrP  = "学号或密码错误 请重试"
	wrongState = "状态异常"

	failLogout        = "登出失败\t%s\n"
	failLogoutBySID   = "SID登出失败\t%s\n"
	failCookieExpired = "Cookie已失效"
	failGetInfo       = "获取已登录账号信息失败"
)
