package tool

import (
	. "ipgw/base"
	"ipgw/tool/get"
	"ipgw/tool/list"
	"ipgw/tool/remove"
	"ipgw/tool/update"
)

var CmdTool = &Command{
	Name:      "tool",
	UsageLine: "ipgw tool",
	Short:     "工具管理",
	Long:      `提供管理工具的功能`,
	Commands: []*Command{
		get.CmdToolGet,
		list.CmdToolList,
		remove.CmdToolRemove,
		update.CmdToolUpdate,
	},
}

/*
	tool get	下载指定的工具
	tool remove	删除指定的工具
	tool list	列出已下载的工具
	tool update	更新已下载的工具

	下载的tool都放在ipgwTool目录下，软(硬)链接到与ipgw同级的目录
	remove时根据ipgwTool目录里的tools.json为根据
	ipgwTool/
		tools.json
		teemo/
		xxx/
*/

// GET逻辑:
// 	获取到tools列表
// 	查找对应tool
//	若存在，校验api版本是否兼容
//	若兼容，查看tool是否已存在./ipgwTool/tools.json
//	若不存在，下载
// 	所有工具的发布为.zip包，下载到./ipgwTool
//	下载后解压到对应目录
//		例如:
//		./ipgwTool/teemo/teemo.exe
// 	然后软连接到./
// 		例如:
// 		./teemo.exe
// 	然后删除压缩包并添加对应的info到./ipgwTool/tools.json

// REMOVE逻辑
//	获取到./ipgwTool/tools.json
//	查看是否存在
//	若存在，删除软连接
//		例如:
//		./teemo.exe
//	再删除本体目录
//		例如
//		./ipgwTool/teemo
//	从./ipgwTool/tools.json中删除对应info

// UPDATE逻辑
//	提示必须先关闭正在运行中的工具（也可以报错之后再输出提示，懒
//	获取到./ipgwTool/tools.json
//	查看是否存在&获取到path
//	获取path下的info.json
//	比对版本是否有更新
//	若有更新，校验api是否兼容
//	若兼容，下载新版本压缩包到./ipgwTool
//	【暂时】先删去旧版本目录再解压缩，这部分同GET【日后可能会有配置文件需求，可能要配合tools.json增加配置文件路径字段，暂时不需要】
//	更新./ipgwTool/tools.json里的info

// LIST逻辑
//	获取到./ipgwTool/tools.json
//	输出info

// LIST -A逻辑
//	获取到tools.json
//	输出info
