package lib

var (
	failDownload = "下载失败，请重试"
	failCreate   = "创建文件失败\t%s\n"
	failConnect  = "网络请求失败\t%s\n"

	fatalLoadToolInfo = "工具列表读取失败"
	fatalSaveToolInfo = "工具列表保存失败"

	wrongUrl = "下载地址错误"

	errNetwork   = "请检查网络连接"
	errEnvReason = "错误原因\t%s\n获取运行环境失败\n"

	failOSNotSupported   = "当前系统暂无发布包"
	failArchNotSupported = "当前架构暂无发布包"
	failQuery            = "获取失败"

	changelog        = "更新日志"
	changelogTitle   = "  %s\n"
	changelogContent = "     %s\n"
)
