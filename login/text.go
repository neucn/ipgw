package login

var (
	mustUsePWhenUseU = "使用账号登陆时请务必使用 -p 输入密码"
	noStoredAccount  = "没有已保存的可用账号"

	usingUP        = "使用账号登陆\t%s\n"
	usingC         = "使用Cookie登陆\t%s\n"
	sendingRequest = "发送登陆请求中..."

	failLogin   = "登陆失败"
	failGetInfo = "用户身份获取失败"
	failOverdue = "余额不足 请充值"

	successGetUsername = "用户身份\t%s\n"
	successLogin       = "登陆成功\t%s\n"
	successGetIP       = "获取IP成功\t%s\n"
	successGetSID      = "获取SID成功\t%s\n"
)
