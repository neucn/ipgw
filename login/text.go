package login

var (
	mustUsePWhenUseU = "使用账号登陆时请务必使用 -p 输入密码"

	tipBeginWithUP = "使用账号登陆\t%s\n"
	tipBeginWithC  = "使用Cookie登陆 %s\n"
	tipRequest     = "发送登陆请求中..."

	errWhenReadLT  = "获取参数时错误: %v\n"
	errWhenRequest = "发送请求时错误: %v\n"
	errNetwork     = "请检查网络连接"

	failGetCookie     = "获取Cookie失败"
	failCookieExpired = "Cookie已失效"

	wrongUOrP = "学号或密码错误，请重试"

	successGetLT       = "获取参数成功\t: %s\n"
	successGetCookie   = "获取Cookie成功\t: %s\n"
	successGetUsername = "用户身份\t: %s\n"
	successLogin       = "登陆成功\t%s\n"
)
