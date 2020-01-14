package get

var (
	infoFetchingList  = "正在获取工具列表..."
	infoSearchingTool = "\n正在搜索工具\t%s\n"
	infoParsingInfo   = "正在解析下载信息"

	infoRequiredAPI = "要求API版本\t%s\n"
	infoLocalAPI    = "本地API版本\t%s\n"

	infoSuccessGet = "下载成功\t%s\n"

	failAlreadyExist     = "本地已有\t%s\n"
	failNoSuchTool       = "无此工具\t%s\n"
	failAPINotCompatible = "API版本不兼容"

	failUnzip   = "解压文件失败\t%s\n"
	failRemove  = "删除压缩包失败\t%s\n"
	failSymlink = "建立软链接失败\t%s\n"
)
