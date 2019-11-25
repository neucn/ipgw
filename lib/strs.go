// 输出文本都在这里统一管理

package lib

const (
	// 1.Error
	// 1.1 登陆
	ErrorCheckConnect  = "请检查网络是否连接"
	ErrorLoginFail     = "登陆失败"
	ErrorNotLoginYet   = "尚未登陆"
	ErrorFetchInfoFail = "信息获取失败"

	// 1.2 配置
	ErrorGetPath        = "在获取环境时出错"
	ErrorLoadUserInfo   = "配置文件读取失败"
	ErrorOpenUserInfo   = "配置文件打开失败"
	ErrorSaveUserInfo   = "在写入配置时出错"
	ErrorUserInfoFormat = "配置文件格式不正确"

	// 1.3 注销
	ErrorLogoutFail = "登出失败"

	// 2.Info
	// 2.1 配置
	InfoLoadUserInfo        = "读取用户配置中..."
	InfoSaveUserInfo        = "存储用户配置中..."
	InfoSaveUserInfoSuccess = "存储用户配置成功"

	// 2.2 流程
	InfoBeginLogin     = "%s, 开始登陆...\n"
	InfoLoginSuccess   = "%s, 登陆成功！\n"
	InfoBeginLogout    = "%s, 开始登出...\n"
	InfoLogoutSuccess  = "%s, 登出成功！\n"
	InfoKickOutSuccess = "%s, 下线成功\n"
	InfoPause          = "请输入回车退出..."

	// 0.Other
	Title     = "IPGW Tool"
	Version   = "Version: 0.0.2"
	Delimiter = ";"
)
