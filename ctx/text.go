package ctx

var (
	fatalGetPath    = "在获取环境时出错"
	fatalLoadInfo   = "配置文件读取失败"
	fatalOpenFile   = "配置文件打开失败"
	fatalSaveInfo   = "在写入配置时出错"
	fatalInfoFormat = "配置文件格式不正确，请尝试使用 ipgw fix 修复"

	errGetCookie = "获取Cookie失败"

	successGetCookie = "获取Cookie成功\t%s\n"

	infoLoading = "读取配置中..."
	infoSaving  = "存储配置中..."
	infoSaved   = "存储配置成功"
)
