package remove

import (
	. "ipgw/base"
	"ipgw/help"
)

var CmdToolRemove = &Command{
	Name:      "remove",
	UsageLine: "ipgw tool remove tool1 tool2 ...",
	Short:     "移除本地工具",
	Long: `提供移除本地工具的功能

  ipgw tool remove teemo
    从本机删除teemo
`,
}

func init() {
	CmdToolRemove.Run = runToolRemove
}

func runToolRemove(cmd *Command, args []string) {
	// 如果没有指定工具名
	if len(args) < 1 {
		help.PrintSimpleUsage(cmd)
	}

	// 处理业务逻辑
	removeAll(args)
}
