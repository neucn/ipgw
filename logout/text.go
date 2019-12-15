package logout

var (
	mustUsePWhenUseU = "使用账号登出时请务必使用 -p 输入密码"

	errWhenReadLT  = "获取参数时错误: %v\n"
	errWhenRequest = "发送请求时错误: %v\n"
	errNetwork     = "请检查网络连接"

	tipBeginWithUP = "使用账号登出\t%s\n"
	tipGetSID      = "获取SID中..."

	successGetLT  = "获取参数成功\t: %s\n"
	successLogout = "登出成功:\t%s\n"

	failLogout        = "登出时失败: %s\n"
	failCookieExpired = "Cookie已失效或不匹配"
)
