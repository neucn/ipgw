package update

var (
	querying = "获取最新版本信息中..."
	forcing  = "强制更新模式"
	updating = "正在更新..."
	removing = "正在移除旧程序..."
	covering = "正在生效新程序..."

	alreadyLatest = "当前已是最新版本"
	localVersion  = "本地版本\t%s\n"
	latestVersion = "最新版本\t%s\n"
	changelog     = "更新日志"

	wrongUrl = "下载地址错误"

	changelogTitle   = "  %s\n"
	changelogContent = "     %s\n"

	errNet    = "请检查网络连接"
	errReason = "错误原因\t%s\n"
	errRunEnv = "获取运行环境失败"

	failQuery            = "获取失败"
	failOSNotSupported   = "当前系统暂无发布包"
	failArchNotSupported = "当前架构暂无发布包"
	failCreate           = "创建文件时失败"
	failUpdate           = "更新失败"
	failDownload         = "下载失败，请重试"
	failChmod            = "执行权限修改失败"

	successUpdate = "更新成功"
)
