package remove

import (
	. "ipgw/base"
	. "ipgw/lib"
	"os"
	"runtime"
)

// 移除所有的工具
func removeAll(toolNames []string) {
	// 获取本地列表
	InfoL(infoFetchingList)
	localList := &Tools{}
	localList.Load()

	// 遍历参数，移除指定的工具
	var tool *Tool
	var ok bool
	for _, name := range toolNames {
		InfoF(infoRemovingTool, name)
		tool, ok = localList.List[name]
		if !ok {
			// 如果不存在，报错并跳到下一个
			ErrorF(failNoSuchTool, name)
			continue
		}

		// 移除工具
		InfoF(infoBeginRemove, name)
		if !removeTool(tool) {
			// 删除失败，跳到下一个
			continue
		}

		// 从localList中删去
		delete(localList.List, name)
		InfoF(infoSuccessRemove, name)
	}

	// 保存localList
	localList.Save()
}

// 移除工具
func removeTool(tool *Tool) (ok bool) {
	// TODO 尽快把info.json里的names的设计给去掉，直接判断如果是windows加个.exe的后缀得了
	// 获取ipgw所在目录
	_, dir := GetRealPathAndDir()
	// 软链接路径
	symlinkPath := dir + tool.Name
	// 判断是否为windows
	if runtime.GOOS == "windows" {
		symlinkPath += ".exe"
	}

	// 删除软链接
	err := os.Remove(symlinkPath)
	if err != nil {
		ErrorF(failRemoveSymlink, err)
		return false
	}

	// 获取ipgwTools目录
	toolsDir := GetToolsDir()

	// 删除工具的整个目录
	err = os.RemoveAll(toolsDir + tool.Name)
	if err != nil {
		ErrorF(failRemoveDir, err)
		return false
	}

	// 删除成功
	return true
}
