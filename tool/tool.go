package tool

import (
	"ipgw/base"
)

func init() {
	CmdTool.Run = runTool
}

var CmdTool = &base.Command{
	Name:      "tool",
	UsageLine: "ipgw tool <command> [argument]",
	Short:     "工具管理",
	Long: `提供管理工具的功能
	tool get	下载指定的工具 
	tool remove	删除指定的工具
	tool list	列出已下载的工具
	tool update	更新已下载的工具
`,
}

func runTool(cmd *base.Command, args []string) {

}
