package cas

var (
	infoSendingLoginRequest = "发送登陆请求中..."

	infoDeviceNotFound = "暂不支持 %s, 自动转换为匿名设备\n"

	errNetwork      = "请检查网络连接"
	errWhenReadArgs = "获取参数时错误\t%v\n"

	failGetArgs = "获取参数失败"
	failGetResp = "获取响应失败"

	failCookieExpired = "Cookie已失效"
	failWrongUOrP     = "学号或密码错误 请重试"
	failWrongSetting  = "一网通设置有误 请打开网页登陆查看"
	failBanned        = "一网通系统报错 可能被ban"

	successGetArgs = "获取参数成功\t%s\n"
)
