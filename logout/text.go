package logout

var (
	mustUsePWhenUseU = "使用账号登出时请务必使用 -p 输入密码"
	noStoredAccount  = "没有已保存的可用账号"

	usingUP  = "使用账号登出\t%s\n"
	usingC   = "使用Cookie登出\t%s\n"
	usingSID = "使用SID登出\t%s\n"

	successLogout      = "登出成功\t%s\n"
	successLogoutBySID = "登出成功"

	failLogout      = "登出失败\t%s\n"
	failLogoutBySID = "SID登出失败\t%s\n"
	failGetInfo     = "获取已登录账号信息失败"
)
