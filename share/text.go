package share

var (
	beginLogout    = "开始登出\t%s\n"
	gettingInfo    = "获取账号信息中..."
	sendingRequest = "发送登陆请求中..."

	errState      = "登陆状态异常"
	errRequest    = "发送请求时错误\t%v\n"
	errNetwork    = "请检查网络连接"
	errWhenReadLT = "获取参数时错误\t%v\n"
	errNoLT       = "无法获取参数"

	differentU     = "不同账号在线\t%s\n"
	deviceNotFound = "暂不支持 %s, 自动转换为匿名设备\n"

	failGetInfo = "获取该账号信息失败"
	failLogout  = "登出失败\t%s\n"
	failGetResp = "获取响应失败"

	wrongUOrP = "学号或密码错误 请重试"
	wrongLock = "一网通设置有误 请打开网页登陆查看"
	wrongBan  = "一网通系统报错 可能被ban"

	successLogout  = "登出成功\t%s\n"
	successGetLT   = "获取参数成功\t%s\n"
	successGetIP   = "获取ip成功\t%s\n"
	successGetSID  = "获取SID成功\t%s\n"
	successGetInfo = "获取信息成功"
)
