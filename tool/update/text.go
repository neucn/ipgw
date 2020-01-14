package update

var (
	infoCloseRunning  = "请关闭正在运行中的工具"
	infoProcessing    = "\n正在处理\t%s\n"
	infoForceUpdate   = "强制更新\t%s\n"
	infoBeginUpdate   = "开始更新\t%s\n"
	infoSuccessUpdate = "更新成功\t%s\n当前版本\t%s\n"

	failNoSuchTool       = "无此工具\t%s\n"
	failNoNewVersion     = "无新版本\t%s\n"
	failAPINotCompatible = "API不兼容\t%s\n"
	failUnzip            = "解压文件失败\t%s\n"
	failRemove           = "删除压缩包失败\t%s\n"
)
