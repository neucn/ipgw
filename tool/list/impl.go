package list

import (
	. "ipgw/base"
	. "ipgw/lib"
)

// 打印本地的工具列表
func printLocalList() {
	// 实例化并载入tools.json
	tools := &Tools{}
	tools.Load()
	// 迭代输出
	InfoL(infoLocalToolList)
	for _, t := range tools.List {
		InfoF(localListText, t.Name, t.Description, t.Version, t.API, t.Author)
	}

}

// 打印在线工具列表，all表示是否输出所有工具
func printOnlineList(all bool) {
	InfoL(infoFetchingList)
	toolList := GetOnlineToolList()
	if all {
		InfoL(infoAllToolList)
		// 若输出全部，则不判断是否兼容
		for _, t := range toolList.List {
			InfoF(onlineListText, t.Name, t.Description, t.API, t.Author)
		}
	} else {
		InfoL(infoAvailableToolList)
		for _, t := range toolList.List {
			// 判断是否兼容，兼容才输出
			if IsAPICompatible(t.API) {
				InfoF(onlineListText, t.Name, t.Description, t.API, t.Author)
			}
		}
	}
}
