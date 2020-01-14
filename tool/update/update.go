package update

import (
	. "ipgw/base"
	"ipgw/help"
)

var CmdToolUpdate = &Command{
	Name:      "update",
	UsageLine: "ipgw tool update [-f force] tool1 tool2 ...",
	Short:     "更新本地工具",
	Long: `提供更新本地工具的功能
  -f    强制更新

  ipgw tool update teemo
    检查更新并升级teemo
  ipgw tool update -f teemo
    强制更新teemo
`,
}

var (
	f bool
)

func init() {
	CmdToolUpdate.Flag.BoolVar(&f, "f", false, "")

	CmdToolUpdate.Run = runToolUpdate
}

func runToolUpdate(cmd *Command, args []string) {
	// 如果没有指定工具名
	if len(args) < 1 {
		help.PrintSimpleUsage(cmd)
	}

	// 处理业务逻辑
	updateAll(args, f)
}
