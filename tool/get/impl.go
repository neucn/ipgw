package get

import (
	. "ipgw/base"
	. "ipgw/lib"
	"os"
	"runtime"
)

// 下载所有的工具
func getAll(toolNames []string) {
	// 先获取工具列表
	InfoL(infoFetchingList)
	localList, onlineList := fetchList()

	// 遍历参数，下载指定的工具
	var tool *Tool
	var ok bool
	for _, name := range toolNames {
		InfoF(infoSearchingTool, name)
		// 检查是否已有该工具
		if _, ok = localList.List[name]; ok {
			ErrorF(failAlreadyExist, name)
			continue
		}

		// 查看该工具是否存在
		tool, ok = onlineList.List[name]
		if !ok {
			// 没有找到则跳过
			ErrorF(failNoSuchTool, name)
			continue
		}

		// 检查API版本是否兼容
		if !IsAPICompatible(tool.API) {
			ErrorL(failAPINotCompatible)
			InfoF(infoRequiredAPI, tool.API)
			InfoF(infoLocalAPI, API)
			continue
		}

		// 检查成功，获取工具
		InfoL(infoParsingInfo)
		if !getTool(tool) {
			// 如果没有获取成功则跳到下一个
			continue
		}
		// 运行到这里说明下载成功，往localList里加信息
		localList.List[tool.Name] = tool
		InfoF(infoSuccessGet, tool.Name)
	}

	// 保存修改后的本地工具列表
	localList.Save()
}

// 获取工具列表
func fetchList() (localList, onlineList *Tools) {
	localList = &Tools{}
	localList.Load()
	onlineList = GetOnlineToolList()
	return
}

// 下载工具
func getTool(tool *Tool) (ok bool) {
	// 获取工具的ver对象
	v := ParseVer(tool.Path)

	// 获取下载地址
	// 有.zip后缀
	url := GetToolDownloadUrl(tool.Path, tool.Name, v)

	// ipgwTool 路径
	toolsDir := GetToolsDir()
	// 下载的临时文件名
	tmpName := toolsDir + tool.Name + ".zip"
	// 目标目录
	targetDir := toolsDir + tool.Name

	// 下载
	Download(url, tmpName)

	// 解压到指定目录
	err := Unzip(tmpName, targetDir)
	if err != nil {
		ErrorF(failUnzip, err)
		return false
	}

	// 解压成功后删除压缩包
	err = os.Remove(tmpName)
	if err != nil {
		ErrorF(failRemove, err)
		return false
	}

	// 下载的工具的可执行文件路径
	executePath := targetDir + string(os.PathSeparator) + v.Name[runtime.GOOS]

	// 赋予可执行文件权限
	_ = os.Chmod(executePath, 0777)

	// 软链接
	// 获取ipgw所在目录
	_, dir := GetRealPathAndDir()
	err = os.Symlink(executePath, dir+v.Name[runtime.GOOS])
	if err != nil {
		ErrorF(failSymlink, err)
		return false
	}

	// 将版本号冗余到tool方便写入tools.json
	tool.Version = v.Latest
	return true
}
