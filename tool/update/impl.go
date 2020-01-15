package update

import (
	. "ipgw/base"
	. "ipgw/lib"
	"os"
	"runtime"
)

// 检查并更新所有的
func updateAll(toolNames []string, force bool) {
	// 提示必须关闭正在运行中的工具
	InfoL(infoCloseRunning)

	// 获取本地tools.json内容
	toolList := &Tools{}
	toolList.Load()

	// 遍历参数，更新
	var tool *Tool
	var ok bool
	for _, name := range toolNames {
		InfoF(infoProcessing, name)
		tool, ok = toolList.List[name]
		if !ok {
			// 若本地不存在该工具信息
			ErrorF(failNoSuchTool, name)
			continue
		}

		// 检查版本
		v := checkUpdate(tool)
		if !v.Update {
			ErrorF(failNoNewVersion, name)
			if !force {
				continue
			}
			InfoF(infoForceUpdate, name)
		}

		// 检查API兼容
		if !IsAPICompatible(v.API) {
			ErrorF(failAPINotCompatible, name)
			continue
		}

		InfoF(infoBeginUpdate, name)
		// 更新
		updateTool(tool, v)
		InfoF(infoSuccessUpdate, name, tool.Version)
	}
	toolList.Save()
}

// 检查是否有最新版本
func checkUpdate(tool *Tool) (v *Ver) {
	// 获取工具的ver对象
	v = ParseVer(tool.Path)

	// 比对版本
	if !IsNewer(v.Latest, tool.Version) {
		// 若并不是最新的
		return
	}
	v.Update = true
	return
}

// 更新
func updateTool(tool *Tool, v *Ver) (ok bool) {
	// 【暂时】覆盖掉tool的整个目录，以后有单独的配置文件需求再说
	// 输出更新日志
	PrintChangelog(tool.Version, v)

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

	// 将版本号冗余到tool方便写入tools.json
	tool.Version = v.Latest
	return true
}
